package gquery

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe/cmd"
	. "github.com/hkoosha/giraffe/core/t11y/dot"
	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/queryerrors"
	"github.com/hkoosha/giraffe/internal/queryimpl"
)

var uintRegex = regexp.MustCompile(`^\d+$`)

type state struct {
	ref     strings.Builder
	flags   cmd.QFlag
	isFin   bool
	noCmd   bool
	escaped bool
	inside  bool
}

type parser struct {
	spec     string
	path     []GiraffeQuery
	state    state
	global   cmd.QFlag
	segment  int
	i        int
	maxDepth uint16
	c        byte
	level    int
}

func (p *parser) onEscape() error {
	switch { //nolint:gocritic
	case p.state.isFin:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)
	}

	p.state.noCmd = true
	p.state.escaped = true

	return nil
}

func (p *parser) onEscaped() error {
	if p.state.isFin || !p.state.noCmd {
		panic(EF("unreachable: invalid escape in query"))
	}

	p.state.escaped = false
	p.state.ref.WriteByte(p.c)

	return nil
}

func (p *parser) onSelf() error {
	switch { //nolint:gocritic
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)
	}

	p.state.flags |= cmd.QModSelf
	p.state.isFin = true
	p.state.noCmd = true

	return nil
}

func (p *parser) onAppend() error {
	switch {
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.state.flags.IsOverwrite(),
		p.state.flags.IsDelete():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	case p.state.flags.IsAppend():
		return queryerrors.DuplicatedCmdError(p.i, p.spec, p.c)
	}

	p.global |= cmd.QModIndeter
	p.state.flags |= cmd.QModAppend
	p.state.isFin = true

	return nil
}

func (p *parser) onDelete() error {
	switch {
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.state.flags.IsOverwrite(),
		p.state.flags.IsAppend(),
		p.state.flags.IsMake(),
		len(p.path) > 0:
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	case p.global.IsDelete():
		return queryerrors.DuplicatedCmdError(p.i, p.spec, p.c)
	}

	p.global |= cmd.QModDelete
	p.state.isFin = true

	return nil
}

func (p *parser) onMake() error {
	switch {
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.state.flags.IsMaybe(),
		p.global.IsDelete():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	case p.state.flags.IsMake():
		return queryerrors.DuplicatedCmdError(p.i, p.spec, p.c)
	}

	p.state.flags |= cmd.QModeMake
	p.state.flags = cmd.QModeMake

	return nil
}

func (p *parser) onMaybe() error {
	switch {
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.state.flags.IsOverwrite(),
		p.state.flags.IsMake():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	case p.state.flags.IsMaybe():
		return queryerrors.DuplicatedCmdError(p.i, p.spec, p.c)
	}

	p.state.flags = cmd.QModeMaybe

	return nil
}

func (p *parser) onOverwrite() error {
	switch {
	case p.state.noCmd:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.global.IsMaybe(),
		p.global.IsAppend(),
		p.global.IsMake(),
		p.global.IsDelete():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)
	}

	p.global = cmd.QModOverwrit

	return nil
}

func (p *parser) onRune() error {
	p.state.ref.WriteByte(p.c)

	return nil
}

func (p *parser) onSep() error {
	str := p.state.ref.String()

	switch {
	case p.state.flags.IsAppend() && str == "":
		str = "0"

	case str == "":
		return queryerrors.EmptyError(p.i, p.spec)
	}

	if int64(p.global.Seq()) >= int64(p.maxDepth) {
		return queryerrors.NestingTooDeepError(p.i, p.spec)
	}

	curr := newQuery(
		nil,
		str,
		p.global|p.state.flags,
	)

	switch {
	case uintRegex.MatchString(str):
		// Not entirely sound, or rather too restrictive if the isMake switch
		// was turned on by previous key parts.
		if !curr.flags.IsMake() && curr.flags.IsMaybe() && curr.ref != "0" {
			return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)
		}

		curr.flags |= cmd.QModArr

		value := cmd.QFlag(M(strconv.ParseUint(str, 10, 64)))
		if value&cmd.ValueMask != value {
			return EF("value too big: %s", str)
		}
		curr.flags |= value

	case p.state.flags.IsAppend():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	default:
		curr.flags |= cmd.QModObj
	}

	if p.global.IsDelete() && p.state.flags.IsMaybe() && !curr.flags.IsArr() {
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)
	}

	p.path = append(p.path, curr)

	//nolint:exhaustruct // is made specifically for zero state.
	p.state = state{}

	seq := p.global.Seq()
	p.global &= ^cmd.SequenceMask
	p.global |= cmd.QFlag((seq + 1) << cmd.SeqShift) //nolint:gosec

	return nil
}

func (p *parser) onBraceR() error {
	switch {
	case !p.state.inside:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.state.ref.Len() == 0:
		return queryerrors.EmptyError(p.i, p.spec)

	case int64(p.global.Seq()) >= int64(p.maxDepth):
		return queryerrors.NestingTooDeepError(p.i, p.spec)
	}

	p.state.inside = false

	subQ, err := mkParser(p.level+1, p.state.ref.String()).parse()
	if err != nil {
		return err
	}

	curr := newQuery(
		nil,
		"["+subQ.String()+"]",
		p.global|p.state.flags|cmd.QModSubQuery,
	)

	p.path = append(p.path, curr)

	//nolint:exhaustruct // is made specifically for zero state.
	p.state = state{}

	seq := p.global.Seq()
	p.global &= ^cmd.SequenceMask
	p.global |= cmd.QFlag((seq + 1) << cmd.SeqShift) //nolint:gosec

	return nil
}

func (p *parser) onBraceL() error {
	switch {
	case p.state.inside:
		// Nested bracket not supported
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	case p.level > 0:
		return queryerrors.NestingTooDeepError(p.i, p.spec)
	}

	p.global |= cmd.QModBraces | cmd.QModDyn
	p.state.inside = true

	return nil
}

func (p *parser) preParse() (bool, error) {
	switch {
	case p.state.escaped:
		if err := p.onEscaped(); err != nil {
			return false, err
		}

	default:
		return false, nil
	}

	return true, nil
}

func (p *parser) doParse() error {
	// for dialect
	p.i++
	p.i += len(dialects.Giraffe1v1.String())

	for p.i = range p.spec {
		p.c = p.spec[p.i]
		switch consumed, err := p.preParse(); {
		case err != nil:
			return err

		case consumed:
			continue

		case p.state.inside && p.c != cmd.BraceR.Byte():
			p.state.ref.WriteByte(p.c)
			continue
		}

		switch p.c {
		case cmd.Escape.Byte():
			if err := p.onEscape(); err != nil {
				return err
			}

		case cmd.Self.Byte():
			if err := p.onSelf(); err != nil {
				return err
			}

		case cmd.Append.Byte():
			if err := p.onAppend(); err != nil {
				return err
			}

		case cmd.Delete.Byte():
			if err := p.onDelete(); err != nil {
				return err
			}

		case cmd.Make.Byte():
			if err := p.onMake(); err != nil {
				return err
			}

		case cmd.Maybe.Byte():
			if err := p.onMaybe(); err != nil {
				return err
			}

		case cmd.Overwrite.Byte():
			if err := p.onOverwrite(); err != nil {
				return err
			}

		case cmd.Sep.Byte():
			if err := p.onSep(); err != nil {
				return err
			}

		case cmd.BraceL.Byte():
			if err := p.onBraceL(); err != nil {
				return err
			}

		case cmd.BraceR.Byte():
			if err := p.onBraceR(); err != nil {
				return err
			}

		case cmd.At.Byte():
			return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

		default:
			if err := p.onRune(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *parser) parsePostValidate() error {
	switch {
	case len(p.path) == 0:
		return queryerrors.EmptyError(p.i, p.spec)

	case p.segment != 0 && p.segment != 2:
		return queryerrors.UnexpectedTokenError(p.i, p.spec, p.c)

	default:
		return nil
	}
}

func (p *parser) postProcess() {
	p.path = slices.Clip(p.path)

	if len(p.path) == 1 {
		p.path[0].flags |= cmd.QModSingle
	}

	isMake := false
	for i := range p.path {
		p.path[i].path = &p.path

		if isMake {
			p.path[i].flags |= cmd.QModeMake
		}
		if p.path[i].flags.IsMake() {
			isMake = true
		}

		if p.global.IsIndeterministic() {
			p.path[i].flags |= cmd.QModIndeter
		}

		if i == 0 {
			p.path[i].flags |= cmd.QModRoot
		}

		if i == len(p.path)-1 {
			p.path[i].flags |= cmd.QModLeaf
		}
	}

	// debugPopulateQueries(p.path)
}

func (p *parser) parse() (GiraffeQuery, error) {
	if err := p.doParse(); err != nil {
		invalid := newQuery(nil, "", cmd.QFlag(0))
		return invalid, err
	}

	if err := p.parsePostValidate(); err != nil {
		invalid := newQuery(nil, "", cmd.QFlag(0))
		return invalid, err
	}

	p.postProcess()

	return p.path[0], nil
}

func mkParser(
	level int,
	spec string,
) *parser {
	if !strings.HasSuffix(spec, cmd.Sep.String()) {
		spec += cmd.Sep.String()
	}

	return &parser{
		maxDepth: queryimpl.MaxDepth,
		spec:     spec,
		global:   cmd.QFlag(0),
		path:     make([]GiraffeQuery, 0, 32),
		segment:  0,
		i:        0,
		c:        0,
		level:    level,
		state: state{
			ref:     strings.Builder{},
			flags:   cmd.Zero,
			isFin:   false,
			noCmd:   false,
			escaped: false,
			inside:  false,
		},
	}
}

func Parse(
	spec string,
) (GiraffeQuery, error) {
	return mkParser(0, spec).parse()
}

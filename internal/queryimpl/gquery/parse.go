package gquery

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe/dialects"
	"github.com/hkoosha/giraffe/internal/queryerrors"
	"github.com/hkoosha/giraffe/cmd"
	"github.com/hkoosha/giraffe/qflag"
	. "github.com/hkoosha/giraffe/t11y/dot"
)

var uintRegex = regexp.MustCompile(`^\d+$`)

type state struct {
	ref     strings.Builder
	flags   qflag.QFlag
	isFin   bool
	noCmd   bool
	escaped bool
}

type parser struct {
	spec     string
	path     []GiraffeQuery
	state    state
	global   qflag.QFlag
	segment  int
	i        int
	maxDepth uint16
	c        byte
}

func (p *parser) reset() {
	//nolint:exhaustruct
	p.state = state{}
	p.state.ref.Grow(64)
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

	p.state.flags |= qflag.QModSelf
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

	p.global |= qflag.QModIndeter
	p.state.flags |= qflag.QModAppend
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

	p.global |= qflag.QModDelete
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

	p.state.flags |= qflag.QModeMake
	p.state.flags = qflag.QModeMake

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

	p.state.flags = qflag.QModeMaybe

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

	p.global = qflag.QModOverwrit

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

		curr.flags |= qflag.QModArr

		value := qflag.QFlag(M(strconv.ParseUint(str, 10, 64)))
		if value&qflag.ValueMask != value {
			panic(EF("value too big: %v", value))
		}

		curr.flags |= value

	case p.state.flags.IsAppend():
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)

	default:
		curr.flags |= qflag.QModObj
	}

	if p.global.IsDelete() && p.state.flags.IsMaybe() && !curr.flags.IsArr() {
		return queryerrors.ConflictingCmdError(p.i, p.spec, p.c)
	}

	p.path = append(p.path, curr)

	p.reset()

	seq := p.global.Seq()
	p.global &= ^qflag.SequenceMask
	p.global |= qflag.QFlag((seq + 1) << qflag.SeqShift) //nolint:gosec

	return nil
}

func (p *parser) onRune() error {
	p.state.ref.WriteByte(p.c)

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
		p.path[0].flags |= qflag.QModSingle
	}

	isMake := false
	for i := range p.path {
		p.path[i].path = &p.path

		if isMake {
			p.path[i].flags |= qflag.QModeMake
		}
		if p.path[i].flags.IsMake() {
			isMake = true
		}

		if p.global.IsIndeterministic() {
			p.path[i].flags |= qflag.QModIndeter
		}

		if i == 0 {
			p.path[i].flags |= qflag.QModRoot
		}

		if i == len(p.path)-1 {
			p.path[i].flags |= qflag.QModLeaf
		}
	}

	// debugPopulateQueries(p.path)
}

func (p *parser) parse() (GiraffeQuery, error) {
	if err := p.doParse(); err != nil {
		invalid := newQuery(nil, "", qflag.QFlag(0))
		return invalid, err
	}

	if err := p.parsePostValidate(); err != nil {
		invalid := newQuery(nil, "", qflag.QFlag(0))
		return invalid, err
	}

	p.postProcess()

	return p.path[0], nil
}

func newParser(
	maxDepth uint16,
	spec string,
) *parser {
	//nolint:exhaustruct
	zeroState := state{}

	p := parser{
		maxDepth: maxDepth,
		spec:     spec,
		state:    zeroState,
		global:   qflag.QFlag(0),
		path:     make([]GiraffeQuery, 0, 32),
		segment:  0,
		i:        0,
		c:        0,
	}
	p.reset()

	return &p
}

func Parse(
	maxDepth uint16,
	spec string,
) (GiraffeQuery, error) {
	if !strings.HasSuffix(spec, cmd.Sep.String()) {
		spec += cmd.Sep.String()
	}
	return newParser(maxDepth, spec).parse()
}

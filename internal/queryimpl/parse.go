package queryimpl

import (
	"slices"
	"strconv"
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot0"
	"github.com/hkoosha/giraffe/internal/queryimpl/gqcache"
	"github.com/hkoosha/giraffe/qcmd"
	"github.com/hkoosha/giraffe/qflag"
)

type state struct {
	ref     strings.Builder
	flags   qflag.QFlag
	isFin   bool
	noCmd   bool
	escaped bool
}

type parser struct {
	spec    string
	path    []Query
	state   state
	global  qflag.QFlag
	segment int

	i int
	c rune

	waitingExternClose bool
}

func (p *parser) reset() {
	//nolint:exhaustruct
	p.state = state{}
	p.state.ref.Grow(64)
}

func (p *parser) last() *Query {
	return &p.path[len(p.path)-1]
}

func (p *parser) newError(
	code uint64,
	msg string,
	extra ...string,
) error {
	sb := strings.Builder{}
	sb.Grow(len(p.spec) + len(msg) + 16)

	sb.WriteString("query parse error, ")
	sb.WriteString(msg)
	sb.WriteString(": at=")
	sb.WriteString(strconv.Itoa(p.i))
	sb.WriteString(", query=")
	sb.WriteString(p.spec)

	for _, e := range extra {
		sb.WriteString(", ")
		sb.WriteString(e)
	}

	return newQueryError(code, sb.String())
}

func (p *parser) unexpectedToken() error {
	return p.newError(
		ErrCodeQueryParseUnexpectedToken,
		"expected token not seen",
		"actual="+string(p.c),
	)
}

func (p *parser) unexpectedSegments() error {
	return p.newError(
		ErrCodeQueryParseUnexpectedSegments,
		"expected number of segments",
		"actual="+strconv.Itoa(p.segment),
	)
}

func (p *parser) unclosedExtern() error {
	return p.newError(
		ErrCodeQueryParseUnclosedExtern,
		"unclosed extern specification",
		"actual="+strconv.Itoa(p.segment),
	)
}

func (p *parser) conflictingCmd(conflictWith rune) error {
	return p.newError(
		ErrCodeQueryParseConflictingCmd,
		"conflicting command",
		"cmd="+string(conflictWith),
	)
}

func (p *parser) duplicatedCmd() error {
	return p.newError(
		ErrCodeQueryParseDuplicatedCmd,
		"duplicated command",
		"cmd="+string(p.c),
	)
}

func (p *parser) emptyQuery() error {
	return p.newError(
		ErrCodeQueryParseEmptyQuery,
		"query is empty",
	)
}

func (p *parser) nestingTooDeep() error {
	return p.newError(
		ErrCodeQueryParseNestingTooDeep,
		"query nesting is too deep",
	)
}

func (p *parser) onEscape() error {
	switch { //nolint:gocritic
	case p.state.isFin:
		return p.unexpectedToken()
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
	p.state.ref.WriteRune(p.c)

	return nil
}

func (p *parser) onSelf() error {
	switch { //nolint:gocritic
	case p.state.noCmd:
		return p.unexpectedToken()
	}

	p.state.flags |= qflag.QModSelf
	p.state.isFin = true
	p.state.noCmd = true

	return nil
}

func (p *parser) onAppend() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(qcmd.Overwrite)

	case p.state.flags.IsAppend():
		return p.duplicatedCmd()

	case p.global.IsDelete():
		return p.conflictingCmd(qcmd.Delete)
	}

	p.global |= qflag.QModIndeter
	p.state.flags |= qflag.QModAppend
	p.state.isFin = true

	return nil
}

func (p *parser) onDelete() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(qcmd.Overwrite)

	case p.state.flags.IsAppend():
		return p.conflictingCmd(qcmd.Delete)

	case p.state.flags.IsMake():
		return p.conflictingCmd(qcmd.Make)

	case p.global.IsDelete():
		return p.duplicatedCmd()

	case len(p.path) > 0:
		return p.conflictingCmd(qcmd.Delete)
	}

	p.global |= qflag.QModDelete
	p.state.isFin = true

	return nil
}

func (p *parser) onMake() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsMaybe():
		return p.conflictingCmd(qcmd.Maybe)

	case p.state.flags.IsMake():
		return p.duplicatedCmd()

	case p.global.IsDelete():
		return p.conflictingCmd(qcmd.Delete)
	}

	p.state.flags |= qflag.QModeMake
	p.state.flags = qflag.QModeMake

	return nil
}

func (p *parser) onMaybe() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(qcmd.Overwrite)

	case p.state.flags.IsMaybe():
		return p.duplicatedCmd()

	case p.state.flags.IsMake():
		return p.conflictingCmd(qcmd.Make)
	}

	p.state.flags = qflag.QModeMaybe

	return nil
}

func (p *parser) onOverwrite() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.global.IsMaybe():
		return p.conflictingCmd(qcmd.Maybe)

	case p.global.IsAppend():
		return p.conflictingCmd(qcmd.Append)

	case p.global.IsMake():
		return p.conflictingCmd(qcmd.Make)

	case p.global.IsDelete():
		return p.conflictingCmd(qcmd.Delete)
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
		return p.emptyQuery()
	}

	if p.global.Seq() >= MaxQueryDepth {
		return p.nestingTooDeep()
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
			return p.conflictingCmd(qcmd.Make)
		}

		curr.flags |= qflag.QModArr

		value := qflag.QFlag(M(strconv.ParseUint(str, 10, 64)))
		if value&qflag.ValueMask != value {
			panic(EF("value too big: %v", value))
		}

		curr.flags |= value

	case p.state.flags.IsAppend():
		return p.conflictingCmd(qcmd.Append)

	default:
		curr.flags |= qflag.QModObj
	}

	if p.global.IsDelete() && p.state.flags.IsMaybe() && !curr.flags.IsArr() {
		return p.conflictingCmd(qcmd.Delete)
	}

	p.path = append(p.path, curr)

	p.reset()

	seq := p.global.Seq()
	p.global &= ^qflag.SequenceMask
	p.global |= qflag.QFlag((seq + 1) << qflag.SeqShift) //nolint:gosec

	return nil
}

func (p *parser) onMove() error {
	p.segment++

	if p.segment > 2 {
		return p.unexpectedToken()
	}

	p.global |= qflag.QModMove

	if err := p.onSep(); err != nil {
		return err
	}

	p.last().flags |= qflag.QModMover

	return nil
}

func (p *parser) onRune() error {
	p.state.ref.WriteRune(p.c)

	return nil
}

func (p *parser) onExtern() error {
	switch {
	case p.waitingExternClose:
		return p.onExternClose()

	default:
		return p.onExternOpen()
	}
}

func (p *parser) onExternOpen() error {
	switch {
	case p.i != 0:
		return p.unexpectedToken()

	case p.state.noCmd:
		return p.unexpectedToken()
	}

	p.waitingExternClose = true
	p.state.flags |= qflag.QModExtern

	return nil
}

func (p *parser) onExternClose() error {
	p.waitingExternClose = false

	return nil
}

func (p *parser) preParse() (bool, error) {
	switch {
	case p.state.escaped:
		if err := p.onEscaped(); err != nil {
			return false, err
		}

	case p.c == qcmd.Extern:
		if err := p.onExtern(); err != nil {
			return false, err
		}

	case p.waitingExternClose:
	// nothing to do

	default:
		return false, nil
	}

	return true, nil
}

func (p *parser) doParse() error {
	for p.i, p.c = range p.spec {
		switch consumed, err := p.preParse(); {
		case err != nil:
			return err

		case consumed:
			continue
		}

		switch p.c {
		case qcmd.Escape:
			if err := p.onEscape(); err != nil {
				return err
			}

		case qcmd.Self:
			if err := p.onSelf(); err != nil {
				return err
			}

		case qcmd.Append:
			if err := p.onAppend(); err != nil {
				return err
			}

		case qcmd.Delete:
			if err := p.onDelete(); err != nil {
				return err
			}

		case qcmd.Make:
			if err := p.onMake(); err != nil {
				return err
			}

		case qcmd.Maybe:
			if err := p.onMaybe(); err != nil {
				return err
			}

		case qcmd.Overwrite:
			if err := p.onOverwrite(); err != nil {
				return err
			}

		case qcmd.Move:
			if err := p.onMove(); err != nil {
				return err
			}

		case qcmd.Sep:
			if err := p.onSep(); err != nil {
				return err
			}

		case qcmd.At:
			return p.unexpectedToken()

		default:
			if err := p.onRune(); err != nil {
				return err
			}
		}
	}

	if p.waitingExternClose {
		return p.unclosedExtern()
	}

	return nil
}

func (p *parser) parsePostValidate() error {
	switch {
	case len(p.path) == 0:
		return p.emptyQuery()

	case p.segment != 0 && p.segment != 2:
		return p.unexpectedSegments()

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
		p.path[i].Path = &p.path

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

	debugPopulateQueries(p.path)
}

func (p *parser) parse() (Query, error) {
	if strings.HasPrefix(p.spec, string(qcmd.Sep)) {
		return ErrQ, p.unexpectedToken()
	}

	if err := p.doParse(); err != nil {
		return ErrQ, err
	}

	if err := p.parsePostValidate(); err != nil {
		return ErrQ, err
	}

	p.postProcess()

	return p.path[0], nil
}

func newParser(spec string) (*parser, error) {
	if !strings.HasSuffix(spec, string(qcmd.Sep)) {
		spec += string(qcmd.At)
	}
	if !strings.ContainsRune(spec, qcmd.Move) {
		spec += string(qcmd.Move) + string(qcmd.At)
	}

	//nolint:exhaustruct
	zeroState := state{}

	p := parser{
		spec:    spec,
		state:   zeroState,
		global:  qflag.QFlag(0),
		path:    make([]Query, 0, 32),
		segment: 0,
		i:       0,
		c:       0,
	}
	p.reset()

	return &p, nil
}

func doParse(
	spec string,
) (Query, error) {
	p, err := newParser(spec)
	if err != nil {
		return ErrQ, err
	}

	q, err := p.parse()
	if err != nil {
		return ErrQ, err
	}

	return q, nil
}

func Parse(
	spec string,
) (Query, error) {
	cached, ok := gqcache.Get[Query](spec)

	if !ok {
		query, err := doParse(spec)
		gqcache.Set(spec, query, err)
		return query, err
	}

	return cached.Unpack()
}

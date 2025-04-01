package gquery

import (
	"slices"
	"strconv"
	"strings"

	. "github.com/hkoosha/giraffe/internal/dot"
)

var commands = map[rune]struct{}{
	CmdOverwrite:        {},
	CmdMake:             {},
	CmdMove:             {},
	CmdMaybe:            {},
	CmdAppend:           {},
	CmdDelete:           {},
	CmdSep:              {},
	CmdEscape:           {},
	CmdAt:               {},
	CmdSelf:             {},
	CmdNonDeterministic: {},
}

type state struct {
	ref     strings.Builder
	flags   QFlag
	isFin   bool
	noCmd   bool
	escaped bool
}

type parser struct {
	spec    string
	path    []Query
	state   state
	global  QFlag
	segment int
	i       int
	c       rune
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

	p.state.flags |= QModSelf
	p.state.isFin = true
	p.state.noCmd = true

	return nil
}

func (p *parser) onAppend() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(CmdOverwrite)

	case p.state.flags.IsAppend():
		return p.duplicatedCmd()

	case p.global.IsDelete():
		return p.conflictingCmd(CmdDelete)
	}

	p.global |= QModNonDet
	p.state.flags |= QModAppend
	p.state.isFin = true

	return nil
}

func (p *parser) onDelete() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(CmdOverwrite)

	case p.state.flags.IsAppend():
		return p.conflictingCmd(CmdDelete)

	case p.state.flags.IsMake():
		return p.conflictingCmd(CmdMake)

	case p.global.IsDelete():
		return p.duplicatedCmd()

	case len(p.path) > 0:
		return p.conflictingCmd(CmdDelete)
	}

	p.global |= QModDelete
	p.state.isFin = true

	return nil
}

func (p *parser) onMake() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsMaybe():
		return p.conflictingCmd(CmdMaybe)

	case p.state.flags.IsMake():
		return p.duplicatedCmd()

	case p.global.IsDelete():
		return p.conflictingCmd(CmdDelete)
	}

	p.state.flags |= QModeMake
	p.state.flags = QModeMake

	return nil
}

func (p *parser) onMaybe() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.state.flags.IsOverwrite():
		return p.conflictingCmd(CmdOverwrite)

	case p.state.flags.IsMaybe():
		return p.duplicatedCmd()

	case p.state.flags.IsMake():
		return p.conflictingCmd(CmdMake)
	}

	p.state.flags = QModeMaybe

	return nil
}

func (p *parser) onOverwrite() error {
	switch {
	case p.state.noCmd:
		return p.unexpectedToken()

	case p.global.IsMaybe():
		return p.conflictingCmd(CmdMaybe)

	case p.global.IsAppend():
		return p.conflictingCmd(CmdAppend)

	case p.global.IsMake():
		return p.conflictingCmd(CmdMake)

	case p.global.IsDelete():
		return p.conflictingCmd(CmdDelete)
	}

	p.global = QModOverwrit

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
			return p.conflictingCmd(CmdMake)
		}

		curr.flags |= QModArr

		value := QFlag(M(strconv.ParseUint(str, 10, 64)))
		if value&ValueMask != value {
			panic(EF("value too big: %v", value))
		}

		curr.flags |= value

	case p.state.flags.IsAppend():
		return p.conflictingCmd(CmdAppend)

	default:
		curr.flags |= QModObj
	}

	if p.global.IsDelete() && p.state.flags.IsMaybe() && !curr.flags.IsArr() {
		return p.conflictingCmd(CmdDelete)
	}

	p.path = append(p.path, curr)

	p.reset()

	seq := p.global.Seq()
	p.global &= ^SequenceMask
	p.global |= QFlag((seq + 1) << seqShift) //nolint:gosec

	return nil
}

func (p *parser) onMove() error {
	p.segment++

	if p.segment > 2 {
		return p.unexpectedToken()
	}

	p.global |= QModMove

	if err := p.onSep(); err != nil {
		return err
	}

	p.last().flags |= QModMover

	return nil
}

func (p *parser) onRune() error {
	p.state.ref.WriteRune(p.c)

	return nil
}

func (p *parser) parse0() error {
	for p.i, p.c = range p.spec {
		if p.state.escaped {
			if err := p.onEscaped(); err != nil {
				return err
			}

			continue
		}

		switch p.c {
		case CmdAt, CmdNonDeterministic:
			return p.unexpectedToken()

		case CmdEscape:
			if err := p.onEscape(); err != nil {
				return err
			}

		case CmdSelf:
			if err := p.onSelf(); err != nil {
				return err
			}

		case CmdAppend:
			if err := p.onAppend(); err != nil {
				return err
			}

		case CmdDelete:
			if err := p.onDelete(); err != nil {
				return err
			}

		case CmdMake:
			if err := p.onMake(); err != nil {
				return err
			}

		case CmdMaybe:
			if err := p.onMaybe(); err != nil {
				return err
			}

		case CmdOverwrite:
			if err := p.onOverwrite(); err != nil {
				return err
			}

		case CmdMove:
			if err := p.onMove(); err != nil {
				return err
			}

		case CmdSep:
			if err := p.onSep(); err != nil {
				return err
			}

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
		p.path[0].flags |= QModSingle
	}

	isMake := false

	for i := range p.path {
		p.path[i].Path = &p.path

		if isMake {
			p.path[i].flags |= QModeMake
		}
		if p.path[i].flags.IsMake() {
			isMake = true
		}

		if p.global.IsNonDeterministic() {
			p.path[i].flags |= QModNonDet
		}

		if i == 0 {
			p.path[i].flags |= QModRoot
		}

		if i == len(p.path)-1 {
			p.path[i].flags |= QModLeaf
		}
	}

	debugPopulateQueries(p.path)
}

func (p *parser) parse() (Query, error) {
	if strings.HasPrefix(p.spec, cmdSepStr) {
		return ErrQ, p.unexpectedToken()
	}

	if err := p.parse0(); err != nil {
		return ErrQ, err
	}

	if err := p.parsePostValidate(); err != nil {
		return ErrQ, err
	}

	p.postProcess()

	return p.path[0], nil
}

func newParser(spec string) (*parser, error) {
	if !strings.HasSuffix(spec, cmdSepStr) {
		spec += cmdSepStr
	}

	//nolint:exhaustruct
	zeroState := state{}

	p := parser{
		spec:    spec,
		state:   zeroState,
		global:  QFlag(0),
		path:    make([]Query, 0, 32),
		segment: 0,
		i:       0,
		c:       0,
	}
	p.reset()

	return &p, nil
}

func parseQueryNoCache(
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

func parse(
	spec string,
) (Query, error) {
	cached, ok := get(spec)

	if !ok {
		query, err := parseQueryNoCache(spec)
		cached = set(spec, query, err)
	}

	return cached.Unpack()
}

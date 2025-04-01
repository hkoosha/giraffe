package gquery

import (
	"strconv"
	"strings"
)

//goland:noinspection SpellCheckingInspection
const (
	QModSelf     QFlag = 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModOverwrit QFlag = 0b01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModNonDet   QFlag = 0b00100000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModeMaybe   QFlag = 0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModeMake    QFlag = 0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModAppend   QFlag = 0b00000100_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModDelete   QFlag = 0b00000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModArr      QFlag = 0b00000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModObj      QFlag = 0b00000000_10000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModRoot     QFlag = 0b00000000_01000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModLeaf     QFlag = 0b00000000_00100000_00000000_00000000_00000000_00000000_00000000_00000000
	QModSingle   QFlag = 0b00000000_00010000_00000000_00000000_00000000_00000000_00000000_00000000
	QModWrite    QFlag = 0b00000000_00001000_00000000_00000000_00000000_00000000_00000000_00000000
	QModMove     QFlag = 0b00000000_00000100_00000000_00000000_00000000_00000000_00000000_00000000
	QModMover    QFlag = 0b00000000_00000010_00000000_00000000_00000000_00000000_00000000_00000000
)

const (
	ValueMask    QFlag = 0b00000000_00000000_00000000_11111111_11111111_11111111_11111111_11111111
	SequenceMask QFlag = 0b00000000_00000000_11111111_00000000_00000000_00000000_00000000_00000000
	ModMask      QFlag = 0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000

	seqShift = 40
)

type QFlag uint64

func (f QFlag) String() string {
	return f.String0()
}

func (f QFlag) String0() string {
	sb := strings.Builder{}

	if f.IsObj() {
		sb.WriteString("{}")
	} else if f.IsArr() {
		sb.WriteByte('[')
		sb.WriteString(strconv.Itoa(f.Val()))
		sb.WriteByte(']')
	} else {
		return "invalid flag: 0b" + strconv.FormatUint(uint64(f), 2)
	}

	f.reconstructPreModIn(&sb)

	sb.WriteByte(CmdAt)
	sb.WriteString(strconv.Itoa(f.Seq()))

	f.reconstructPostModIn(&sb)

	return sb.String()
}

func (f QFlag) reconstructPreMod() string {
	sb := strings.Builder{}
	f.reconstructPreModIn(&sb)

	return sb.String()
}

func (f QFlag) reconstructPreModIn(
	sb *strings.Builder,
) {
	if f&QModOverwrit != 0 {
		sb.WriteByte(CmdOverwrite)
	}

	if f&QModeMaybe != 0 {
		sb.WriteByte(CmdMaybe)
	}

	if f&QModAppend != 0 {
		sb.WriteByte(CmdAppend)
	}

	if f&QModDelete != 0 {
		sb.WriteByte(CmdDelete)
	}

	if f&QModeMake != 0 {
		sb.WriteByte(CmdMake)
	}

	if f&QModSelf != 0 {
		sb.WriteByte(CmdSelf)
	}
}

func (f QFlag) reconstructPostMod() string {
	sb := strings.Builder{}
	f.reconstructPostModIn(&sb)

	return sb.String()
}

func (f QFlag) reconstructPostModIn(
	sb *strings.Builder,
) {
	if f&QModMover != 0 {
		sb.WriteByte(CmdMove)
	}
}

func (f QFlag) IsNonDeterministic() bool {
	return f&QModNonDet != 0
}

func (f QFlag) IsMaybe() bool {
	return f&QModeMaybe != 0
}

func (f QFlag) IsMake() bool {
	return f&QModeMake != 0
}

func (f QFlag) IsAppend() bool {
	return f&QModAppend != 0
}

func (f QFlag) IsDelete() bool {
	return f&QModDelete != 0
}

func (f QFlag) IsOverwrite() bool {
	return f&QModOverwrit != 0
}

func (f QFlag) IsArr() bool {
	return f&QModArr != 0
}

func (f QFlag) IsObj() bool {
	return f&QModObj != 0
}

func (f QFlag) IsSelf() bool {
	return f&QModSelf != 0
}

func (f QFlag) IsRoot() bool {
	return f&QModRoot != 0
}

func (f QFlag) IsLeaf() bool {
	return f&QModLeaf != 0
}

func (f QFlag) IsSingle() bool {
	return f&QModSingle != 0
}

func (f QFlag) IsWrite() bool {
	return f&QModWrite != 0
}

func (f QFlag) IsMove() bool {
	return f&QModMove != 0
}

func (f QFlag) IsMover() bool {
	return f&QModMover != 0
}

func (f QFlag) IsReadonly() bool {
	return !f.IsWrite() && !f.IsMove()
}

func (f QFlag) Val() int {
	//nolint:gosec
	return int(f & ValueMask)
}

func (f QFlag) Seq() int {
	//nolint:gosec
	return int((f & SequenceMask) >> seqShift)
}

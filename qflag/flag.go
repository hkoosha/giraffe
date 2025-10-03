package qflag

import (
	"strconv"
	"strings"

	"github.com/hkoosha/giraffe/qcmd"
)

//goland:noinspection SpellCheckingInspection
const (
	QModSelf     QFlag = 0b10000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModOverwrit QFlag = 0b01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	QModIndeter  QFlag = 0b00100000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
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
	QModExtern   QFlag = 0b00000000_00000001_00000000_00000000_00000000_00000000_00000000_00000000
)

const (
	ValueMask    QFlag = 0b00000000_00000000_00000000_11111111_11111111_11111111_11111111_11111111
	SequenceMask QFlag = 0b00000000_00000000_11111111_00000000_00000000_00000000_00000000_00000000
	ModMask      QFlag = 0b11111111_11111111_00000000_00000000_00000000_00000000_00000000_00000000

	SeqShift = 40
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

	f.ReconstructPreModIn(&sb)

	sb.WriteByte(qcmd.At.Byte())
	sb.WriteString(strconv.Itoa(f.Seq()))

	return sb.String()
}

func (f QFlag) ReconstructPreMod() string {
	sb := strings.Builder{}
	f.ReconstructPreModIn(&sb)

	return sb.String()
}

func (f QFlag) ReconstructPreModIn(
	sb *strings.Builder,
) {
	if f&QModOverwrit != 0 {
		sb.WriteByte(qcmd.Overwrite.Byte())
	}

	if f&QModeMaybe != 0 {
		sb.WriteByte(qcmd.Maybe.Byte())
	}

	if f&QModAppend != 0 {
		sb.WriteByte(qcmd.Append.Byte())
	}

	if f&QModDelete != 0 {
		sb.WriteByte(qcmd.Delete.Byte())
	}

	if f&QModeMake != 0 {
		sb.WriteByte(qcmd.Make.Byte())
	}

	if f&QModSelf != 0 {
		sb.WriteByte(qcmd.Self.Byte())
	}
}

func (f QFlag) IsIndeterministic() bool {
	return f&QModIndeter != 0
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

func (f QFlag) IsExtern() bool {
	return f&QModExtern != 0
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
	return int((f & SequenceMask) >> SeqShift)
}

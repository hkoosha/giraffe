package gdatum

import (
	"strconv"
	"sync/atomic"

	"github.com/hkoosha/giraffe/core/serdes/gson"
)

const debugNilId = "<nil>"

var debugDatumId = atomic.Uint64{}

func init() {
	debugDatumId.Add(11)
}

func NewDatumDebug() DatumDebug {
	return DatumDebug{
		ID: debugDatumId.Add(1),
	}
}

type DatumDebug struct {
	ID uint64 `json:"id"`
}

func (d *DatumDebug) String() string {
	if d == nil {
		return debugNilId
	}

	serialized, err := gson.Marshal(d)
	if err != nil {
		panic(err)
	}

	return strconv.FormatUint(d.ID, 10) + "#" + string(serialized)
}

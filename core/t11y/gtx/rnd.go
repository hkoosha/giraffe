package gtx

import (
	randc "crypto/rand"
	"math/big"
	rand1 "math/rand"
	rand2 "math/rand/v2"

	. "github.com/hkoosha/giraffe/core/t11y/dot"
)

var int63 = big.NewInt(int64(1) << 62)

func seed() (int64, error) {
	r, err := randc.Int(randc.Reader, int63)
	if err != nil {
		return 0, err
	}

	return r.Int64(), nil
}

type rnd struct {
	seed func() (int64, error)
}

func (r rnd) StdV1() *rand1.Rand {
	//nolint:gosec
	return rand1.New(rand1.NewSource(M(seed())))
}

func (r rnd) StdV2() *rand2.Rand {
	//nolint:gosec
	return rand2.New(rand2.NewPCG(
		uint64(M(r.seed())),
		uint64(M(r.seed())),
	))
}

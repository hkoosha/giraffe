package gquery_test

import (
	"math/bits"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe/cmd"
)

func TestQFlag_Mods(t *testing.T) {
	qFlags := []cmd.QFlag{
		cmd.QModIndeter,
		cmd.QModeMaybe,
		cmd.QModeMake,
		cmd.QModAppend,
		cmd.QModArr,
		cmd.QModObj,
	}

	t.Run("bit count", func(t *testing.T) {
		for _, qFlag := range qFlags {
			assert.Equal(
				t,
				1,
				bits.OnesCount64(uint64(qFlag)),
				"%064s", strconv.FormatUint(uint64(qFlag), 2),
			)
		}
	})

	t.Run("mask", func(t *testing.T) {
		for _, qFlag := range qFlags {
			assert.Equal(
				t,
				qFlag,
				qFlag&cmd.ModMask,
			)
		}
	})
}

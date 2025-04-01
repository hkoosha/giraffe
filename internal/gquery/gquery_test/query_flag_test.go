package gquery_test

import (
	"math/bits"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hkoosha/giraffe/internal/gquery"
)

func TestQFlag_Mods(t *testing.T) {
	qFlags := []gquery.QFlag{
		gquery.QModNonDet,
		gquery.QModeMaybe,
		gquery.QModeMake,
		gquery.QModAppend,
		gquery.QModArr,
		gquery.QModObj,
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
				qFlag&gquery.ModMask,
			)
		}
	})
}

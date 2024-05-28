package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkBy(t *testing.T) {
	t.Run("strings chunk by 10", func(t *testing.T) {
		srcArr := make([]string, 0, 20)
		for i := 0; i < 20; i++ {
			srcArr = append(srcArr, strconv.Itoa(i))
		}
		want := [][]string{
			{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
			{"10", "11", "12", "13", "14", "15", "16", "17", "18", "19"},
		}
		got := ChunkBy(srcArr, 10)

		assert.Equal(t, want, got)
	})
}

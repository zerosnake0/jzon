package jzon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Val_ReadVal(t *testing.T) {
	t.Run("nil pointer receiver error", func(t *testing.T) {
		it := NewIterator()
		err := it.ReadVal(nil)
		require.Equal(t, NilPointerReceiverError, err)
	})
	t.Run("pointer receiver error", func(t *testing.T) {
		it := NewIterator()
		var o string
		err := it.ReadVal(o)
		require.Equal(t, PointerReceiverError, err)
	})
	t.Run("struct", func(t *testing.T) {
		it := NewIterator()
		var p struct {
			K string `json:"k"`
		}
		it.ResetBytes([]byte(` { "k": "v" } `))
		err := it.ReadVal(&p)
		require.NoError(t, err)
		require.Equal(t, "v", p.K)
	})
}

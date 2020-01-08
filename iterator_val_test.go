package jzon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Val_ReadVal(t *testing.T) {
	t.Run("nil pointer receiver error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			err := it.ReadVal(nil)
			require.Equal(t, NilPointerReceiverError, err)
		})
	})
	t.Run("pointer receiver error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			var o string
			err := it.ReadVal(o)
			require.Equal(t, PointerReceiverError, err)
		})
	})
	t.Run("struct", func(t *testing.T) {
		withIterator(` { "k": "v" } `, func(it *Iterator) {
			var p struct {
				K string `json:"k"`
			}
			err := it.ReadVal(&p)
			require.NoError(t, err)
			require.Equal(t, "v", p.K)
		})
	})
}

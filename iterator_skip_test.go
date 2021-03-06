package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Skip_Skip(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			err := it.Skip()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("skip", func(t *testing.T) {
		s := `{
			"string": "string",
			"null": null,
			"true": true,
			"false": false,
			"number": -123.0456E+789,
			"array": [ "string", null, true, false,
				-123.0456E+789, [ ], { } ],
			"object": {
				"string": "string",
				"null": null,
				"true": true,
				"false": false,
				"number": -123.0456E789,
				"array": [ "string", null, true, false,
					-123.0456E+789, [ ], { } ],
				"object": {	}
			}
		}`
		withIterator(s, func(it *Iterator) {
			err := it.Skip()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
}

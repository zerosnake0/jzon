package jzon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := stackPool.Get().(*stack).init()
		require.Equal(t, stackElementNone, s.top())
		require.Equal(t, stackElementNone, s.pop())
	})
	t.Run("array top", func(t *testing.T) {
		s := stackPool.Get().(*stack).initArray()
		require.Equal(t, stackElementArrayBegin, s.top())
		require.Equal(t, stackElementArray, s.pop())
		require.Equal(t, stackElementNone, s.top())
		require.Equal(t, stackElementNone, s.pop())
	})
	t.Run("object top", func(t *testing.T) {
		s := stackPool.Get().(*stack).initObject()
		require.Equal(t, stackElementObjectBegin, s.top())
		require.Equal(t, stackElementObject, s.pop())
		require.Equal(t, stackElementNone, s.top())
		require.Equal(t, stackElementNone, s.pop())
	})
	t.Run("nested 1", func(t *testing.T) {
		s := stackPool.Get().(*stack)
		count := 100
		for i := 0; i < count; i++ {
			if i&1 == 0 {
				s.pushObject()
				require.Equal(t, stackElementObjectBegin, s.top())
			} else {
				s.pushArray()
				require.Equal(t, stackElementArrayBegin, s.top())
			}
		}
		for i := count - 1; i >= 0; i-- {
			if i&1 == 0 {
				require.Equal(t, stackElementObjectBegin, s.top())
				require.Equal(t, stackElementObject, s.pop())
			} else {
				require.Equal(t, stackElementArrayBegin, s.top())
				require.Equal(t, stackElementArray, s.pop())
			}
		}
		require.Equal(t, stackElementNone, s.top())
		require.Equal(t, stackElementNone, s.pop())
	})
}

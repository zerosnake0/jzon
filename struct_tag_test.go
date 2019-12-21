package jzon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructTag_parseTag(t *testing.T) {
	t.Run("no comma", func(t *testing.T) {
		name, opts := parseTag("test")
		require.Equal(t, "test", name)
		require.Empty(t, opts)
	})
	t.Run("with comma", func(t *testing.T) {
		name, opts := parseTag("test,opts")
		require.Equal(t, "test", name)
		require.Equal(t, tagOptions("opts"), opts)
	})
}

func TestStructTag_TagOptions_Contains(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var opts tagOptions
		require.False(t, opts.Contains("opt"))
	})
	t.Run("non empty", func(t *testing.T) {
		opts := tagOptions("opt1,opt2,opt3")
		require.True(t, opts.Contains("opt1"))
		require.True(t, opts.Contains("opt2"))
		require.True(t, opts.Contains("opt3"))
		require.False(t, opts.Contains("opt4"))
	})
}

func TestStructTag_isValidTag(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.False(t, isValidTag(""))
	})
	t.Run("true", func(t *testing.T) {
		require.True(t, isValidTag("!#$%&()*+-./:<=>?@[]^_{|}~ "))
	})
	t.Run("false", func(t *testing.T) {
		require.False(t, isValidTag("\u00b6"))
	})
}

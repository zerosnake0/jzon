package jzon

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamer_Float32_Error(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.Float32(1)
		})
	})
	t.Run("infinity", func(t *testing.T) {
		f := float32(math.Inf(0))
		s := NewStreamer()
		defer ReturnStreamer(s)
		s.Float32(f)
		require.Equal(t, FloatIsInfinity, s.Error)
	})
	t.Run("nan", func(t *testing.T) {
		f := float32(math.NaN())
		s := NewStreamer()
		defer ReturnStreamer(s)
		s.Float32(f)
		require.Equal(t, FloatIsNan, s.Error)
	})
}

func TestStreamer_Float32(t *testing.T) {
	f := func(t *testing.T, f float32) {
		checkEncodeWithStandard(t, DefaultEncoder, f, func(s *Streamer) {
			s.Float32(f)
		}, nil)
	}
	t.Run("1.2e-3", func(t *testing.T) {
		f(t, 1.2e-3)
	})
	t.Run("1e-7", func(t *testing.T) {
		f(t, 1e-7)
	})
	t.Run("1e21", func(t *testing.T) {
		f(t, 1e21)
	})
	t.Run("-1e-7", func(t *testing.T) {
		f(t, -1e-7)
	})
	t.Run("-1e21", func(t *testing.T) {
		f(t, -1e21)
	})
}

func TestStreamer_Float64_Error(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.Float64(1)
		})
	})
	t.Run("infinity", func(t *testing.T) {
		f := math.Inf(0)
		s := NewStreamer()
		defer ReturnStreamer(s)
		s.Float64(f)
		require.Equal(t, FloatIsInfinity, s.Error)
	})
	t.Run("nan", func(t *testing.T) {
		f := math.NaN()
		s := NewStreamer()
		defer ReturnStreamer(s)
		s.Float64(f)
		require.Equal(t, FloatIsNan, s.Error)
	})
}

func TestStreamer_Float64(t *testing.T) {
	f := func(t *testing.T, f float64) {
		checkEncodeWithStandard(t, DefaultEncoder, f, func(s *Streamer) {
			s.Float64(f)
		}, nil)
	}
	t.Run("1.2e-3", func(t *testing.T) {
		f(t, 1.2e-3)
	})
	t.Run("1e-7", func(t *testing.T) {
		f(t, 1e-7)
	})
	t.Run("1e21", func(t *testing.T) {
		f(t, 1e21)
	})
	t.Run("-1e-7", func(t *testing.T) {
		f(t, -1e-7)
	})
	t.Run("-1e21", func(t *testing.T) {
		f(t, -1e21)
	})
}

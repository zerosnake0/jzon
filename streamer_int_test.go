package jzon

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamer_Int8(t *testing.T) {
	f := func(t *testing.T, i int8) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Int8(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxInt8)
	})
	t.Run("min", func(t *testing.T) {
		f(t, math.MinInt8)
	})
}

func TestStreamer_Uint8(t *testing.T) {
	f := func(t *testing.T, i uint8) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Uint8(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxUint8)
	})
	t.Run("min", func(t *testing.T) {
		f(t, 0)
	})
}

func TestStreamer_Int16(t *testing.T) {
	f := func(t *testing.T, i int16) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Int16(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxInt16)
	})
	t.Run("min", func(t *testing.T) {
		f(t, math.MinInt16)
	})
}

func TestStreamer_Uint16(t *testing.T) {
	f := func(t *testing.T, i uint16) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Uint16(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxUint16)
	})
	t.Run("min", func(t *testing.T) {
		f(t, 0)
	})
}

func TestStreamer_Int32(t *testing.T) {
	f := func(t *testing.T, i int32) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Int32(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxInt32)
	})
	t.Run("min", func(t *testing.T) {
		f(t, math.MinInt32)
	})
}

func TestStreamer_Uint32(t *testing.T) {
	f := func(t *testing.T, i uint32) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Uint32(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxUint32)
	})
	t.Run("min", func(t *testing.T) {
		f(t, 0)
	})
	t.Run("1023", func(t *testing.T) {
		f(t, 1023)
	})
	t.Run("1023045", func(t *testing.T) {
		f(t, 1023045)
	})
}

func TestStreamer_Int64(t *testing.T) {
	f := func(t *testing.T, i int64) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Int64(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxInt64)
	})
	t.Run("min", func(t *testing.T) {
		f(t, math.MinInt64)
	})
}

func TestStreamer_Uint64(t *testing.T) {
	f := func(t *testing.T, i uint64) {
		exp, err := json.Marshal(i)
		require.NoError(t, err)
		testStreamer(t, string(exp), func(s *Streamer) {
			s.Uint64(i)
		})
	}
	t.Run("max", func(t *testing.T) {
		f(t, math.MaxUint64)
	})
	t.Run("min", func(t *testing.T) {
		f(t, 0)
	})
	t.Run("1023", func(t *testing.T) {
		f(t, 1023)
	})
	t.Run("1023045", func(t *testing.T) {
		f(t, 1023045)
	})
	t.Run("1023045067", func(t *testing.T) {
		f(t, 1023045067)
	})
	t.Run("1023045067089", func(t *testing.T) {
		f(t, 1023045067089)
	})
	t.Run("1023045067089000", func(t *testing.T) {
		f(t, 1023045067089000)
	})
}

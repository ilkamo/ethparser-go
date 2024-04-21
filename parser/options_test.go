package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ilkamo/ethparser-go/tests/mock"
)

func TestWithBlockProcessTimeout(t *testing.T) {
	t.Run("set block process timeout opt", func(t *testing.T) {
		p := NewParser(nil, nil, 0,
			nil, nil, WithBlockProcessTimeout(2000))
		require.Equal(t, time.Duration(2000), p.blockProcessTimeout)
	})
}

func TestWithLogger(t *testing.T) {
	t.Run("set logger opt", func(t *testing.T) {
		log := &mock.Logger{}

		p := NewParser(nil, nil, 0,
			nil, nil, WithLogger(log))
		require.NotNil(t, p.logger)
		require.Equal(t, log, p.logger)
	})
}

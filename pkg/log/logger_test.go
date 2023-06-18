package log

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(DevLogger, zap.DebugLevel)
	require.NoError(t, err)

	logger.Info("Test event")
	logger.Debug("this is a test DEBUG log message.")
}

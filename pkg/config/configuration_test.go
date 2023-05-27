package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewConfig(t *testing.T) {
	conf, err := NewConfig()
	require.NoError(t, err)
	require.NotNil(t, conf.Client)
	require.Equal(t, 30, conf.Client.TimeOut, "wrong timeout value")
	t.Log(conf.Client.TimeOut)

}

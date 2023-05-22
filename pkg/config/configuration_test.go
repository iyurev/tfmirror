package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	conf, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf.Client.TimeOut)

}

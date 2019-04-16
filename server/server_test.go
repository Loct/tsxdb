package server_test

import (
	"../server"
	"testing"
)

func TestNew(t *testing.T) {
	opts := server.NewOpts()
	s := server.New(opts)
	if s == nil {
		t.Error()
	}
	if err := s.Init(); err != nil {
		t.Error(err)
	}
}

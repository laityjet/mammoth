package pctx

import (
	"testing"

	"github.com/laityjet/mammoth/v0/internal/log"
)

func TestBackground(t *testing.T) {
	_, h := log.TestWithCapture(t)
	log.Info(Background(""), "hi")
	h.HasALog(t)
}

func TestTODO(t *testing.T) {
	_, h := log.TestWithCapture(t)
	log.Info(TODO(), "hi")
	h.HasALog(t)
}

package ppsutil_test

import (
	"io"
	"strings"
	"testing"

	"github.com/laityjet/mammoth/v0/internal/errors"
	"github.com/laityjet/mammoth/v0/internal/ppsutil"
	"github.com/laityjet/mammoth/v0/internal/pfs"
	"github.com/laityjet/mammoth/v0/internal/pps"
)

func TestPipelineRcName(t *testing.T) {
	for _, c := range []struct {
		projectName, pipelineName string
		version                   uint64
		name                      string
	}{
		{"default", "foo", 1, "default-foo-v1"},
		{"foo", "bar", 1, "foo-bar-v1"},
	} {
		p := &pps.PipelineInfo{Pipeline: &pps.Pipeline{Project: &pfs.Project{Name: c.projectName}, Name: c.pipelineName}, Version: c.version}
		name := ppsutil.PipelineRcName(p)
		if name != c.name {
			t.Errorf("case %v: expected %q; got %q", c, c.name, name)
		}
	}
}

func Test_PipelineManifestReader(t *testing.T) {
	// NOTE: the spec below uses a string value for the parallelism spec
	// constant.  This follows the Protobuf JSON mapping spec[1], which
	// indicates that for uint64 the “JSON value will be a decimal string.
	// Either numbers or strings are accepted.”  This test ensures that
	// decoding properly handles this mapping.
	//
	// [1] https://protobuf.dev/programming-guides/proto3/#json
	r, err := ppsutil.NewPipelineManifestReader(strings.NewReader(`{
  "pipeline": {
    "name": "first"
  },
  "input": {
    "pfs": {
      "glob": "/*",
      "repo": "input"
    }
  },
  "parallelism_spec": {
    "constant": "1"
  },
  "transform": {
    "cmd": [ "/bin/bash" ],
    "stdin": [
      "cp /pfs/input/* /pfs/out"
    ]
  }
}
`))
	if err != nil {
		t.Error(err)
	}
	var i int
	for {
		p, err := r.NextCreatePipelineRequest()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		if expected, got := uint64(1), p.ParallelismSpec.Constant; expected != got {
			t.Errorf("parallelism spec constant: expected %d; got %d", expected, got)
		}
		i++
	}
	if expected, got := 1, i; expected != got {
		t.Errorf("expected %d objects; got %d", expected, got)
	}
}

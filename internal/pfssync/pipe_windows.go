//go:build windows
// +build windows

package pfssync

import (
	"io"

	"github.com/laityjet/mammoth/v0/internal/errors"
)

func (d *downloader) makePipe(path string, cb func(io.Writer) error) error {
	return errors.Errorf("lazy file download through pipes is not supported on Windows")
}

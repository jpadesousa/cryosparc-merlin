package installer

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func (i *Installer) createHomeDirectory(_ func(float64)) error {
	time.Sleep(500 * time.Millisecond)

	info, err := os.Stat(i.cfg.HomeDir)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return os.MkdirAll(i.cfg.HomeDir, 0o755)

	case err != nil:
		return fmt.Errorf("check home directory: %w", err)

	case !info.IsDir():
		return fmt.Errorf("%q exists but is not a directory", i.cfg.HomeDir)

	default:
		return fmt.Errorf("directory %q already exists", i.cfg.HomeDir)
	}
}

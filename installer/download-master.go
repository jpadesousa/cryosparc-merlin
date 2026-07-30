package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (i *Installer) downloadMaster(update func(float64)) error {
	url := fmt.Sprintf(
		"https://get.cryosparc.com/download/master-v%s/%s",
		i.cfg.Release,
		i.cfg.License,
	)

	archive := filepath.Join(i.cfg.HomeDir, "cryosparc_master.tar.gz")

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download CryoSPARC master package: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"download failed (HTTP %d). Please check the CryoSPARC release, license ID, and internet connection",
			resp.StatusCode,
		)
	}

	f, err := os.Create(archive)
	if err != nil {
		return fmt.Errorf("create %q: %w", archive, err)
	}
	defer f.Close()

	reader := &progressReader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
		Update: update,
	}

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("save %q: %w", archive, err)
	}

	return nil
}

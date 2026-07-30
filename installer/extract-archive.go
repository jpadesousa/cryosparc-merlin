package installer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (i *Installer) extractArchive(_ func(float64), filename string) error {
	time.Sleep(500 * time.Millisecond)

	archive := filepath.Join(
		i.cfg.HomeDir,
		filename,
	)

	file, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("open archive %q: %w", archive, err)
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read gzip archive: %w", err)
	}
	defer gz.Close()

	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return fmt.Errorf("read archive: %w", err)
		}

		target := filepath.Join(
			i.cfg.HomeDir,
			header.Name,
		)

		switch header.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf(
					"create directory %q: %w",
					target,
					err,
				)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf(
					"create parent directory %q: %w",
					target,
					err,
				)
			}

			out, err := os.OpenFile(
				target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				os.FileMode(header.Mode),
			)
			if err != nil {
				return fmt.Errorf(
					"create file %q: %w",
					target,
					err,
				)
			}

			_, copyErr := io.Copy(out, tarReader)
			closeErr := out.Close()

			if copyErr != nil {
				return fmt.Errorf(
					"extract %q: %w",
					target,
					copyErr,
				)
			}

			if closeErr != nil {
				return fmt.Errorf(
					"close %q: %w",
					target,
					closeErr,
				)
			}
		}
	}
}

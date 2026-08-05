package command

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var ErrInstallationCancelled = errors.New("Installation cancelled by user")

// =============================================================================
// Cryosparcm command
// =============================================================================
func cryosparcmCmd(remotehost, homedir string, args ...string) *exec.Cmd {
	remoteCmd := append([]string{"cryosparcm"}, args...)

	if homedir != "" {
		remoteCmd = append(
			[]string{"cryosparcm", "--dir", filepath.Join(homedir, "cryosparc_master")},
			args...,
		)
	}

	cmd := exec.Command(
		"ssh",
		fmt.Sprintf("%s@%s", os.Getenv("USER"), remotehost),
		strings.Join(remoteCmd, " "),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// =============================================================================
// Download Master
// =============================================================================
type progressReader struct {
	io.Reader

	Total   int64
	Current int64
	Update  func(float64)
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)

	if n > 0 {
		r.Current += int64(n)

		if r.Total > 0 {
			r.Update(float64(r.Current) / float64(r.Total))
		}
	}

	return n, err
}

func (i *Command) downloadCryosparc(update func(float64), downloaddir, release, license, installation, arch string) error {
	url := fmt.Sprintf(
		"https://get.cryosparc.com/download/%s-v%s/%s?arch=%s",
		installation,
		release,
		license,
		arch,
	)

	archive := filepath.Join(downloaddir, fmt.Sprintf("cryosparc_%s.tar.gz", installation))

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download CryoSPARC %s package: %w", installation, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"download failed (HTTP %d). Please check the CryoSPARC release, license ID, architecture, and internet connection",
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

// =============================================================================
// Directory Check
// =============================================================================
func directoryExistsAndNotEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil

	case err != nil:
		return false, fmt.Errorf("stat %q: %w", path, err)

	case !info.IsDir():
		return false, fmt.Errorf("%q exists but is not a directory", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	switch {
	case errors.Is(err, io.EOF):
		return false, nil // directory exists but is empty

	case err != nil:
		return false, fmt.Errorf("read %q: %w", path, err)

	default:
		return true, nil // directory exists and contains at least one entry
	}
}

// =============================================================================
// Create Directory
// =============================================================================
func (i *Command) createDirectory(_ func(float64), path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	return nil
}

// =============================================================================
// Extract archive
// =============================================================================
func (i *Command) extractArchive(_ func(float64), downloaddir, extractdir, filename string) error {

	archive := filepath.Join(
		downloaddir,
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
			extractdir,
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

// =============================================================================
// Install Master
// =============================================================================
func (i *Command) installMaster(
	homedir,
	license,
	remotehost,
	dbdir string,
	baseport uint) *exec.Cmd {

	script := filepath.Join(
		homedir,
		"cryosparc_master",
		"install.sh",
	)

	cmd := exec.Command(
		script,
		"--license", license,
		"--hostname", remotehost,
		"--dbpath", dbdir,
		"--port", strconv.Itoa(int(baseport)),
	)

	cmd.Dir = filepath.Dir(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// =============================================================================
// Install Worker
// =============================================================================
func (i *Command) installWorker(homedir, license string) *exec.Cmd {

	script := filepath.Join(
		homedir,
		"cryosparc_worker",
		"install.sh",
	)

	cmd := exec.Command(
		script,
		"--license", license,
	)

	cmd.Dir = filepath.Dir(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

// =============================================================================
// Replace CryoSPARC license in config.sh
// =============================================================================

func (i *Command) replaceLicenseID(_ func(float64), installDir, license string) error {

	path := filepath.Join(installDir, "config.sh")

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config.sh %q: %w", path, err)
	}

	re := regexp.MustCompile(`(?m)^export CRYOSPARC_LICENSE_ID=".*"$`)

	updated := re.ReplaceAllString(
		string(data),
		fmt.Sprintf(`export CRYOSPARC_LICENSE_ID="%s"`, license),
	)

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write config.sh %q: %w", path, err)
	}

	return nil
}

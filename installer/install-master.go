package installer

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

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

func (i *Installer) installMaster() *exec.Cmd {
	script := filepath.Join(
		i.cfg.HomeDir,
		"cryosparc_master",
		"install.sh",
	)

	cmd := exec.Command(
		script,
		"--license", i.cfg.License,
		"--hostname", i.cfg.RemoteHost,
		"--dbpath", i.cfg.DbDir,
		"--port", strconv.Itoa(int(i.cfg.BasePort)),
	)

	cmd.Dir = filepath.Dir(script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

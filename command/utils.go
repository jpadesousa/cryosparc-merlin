package command

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrCancelled = errors.New("cancelled by user")

// =============================================================================
// Cryosparcm command
// =============================================================================
func cryosparcmCmd(hostname, cryosparcpath, commandpath string, help bool, args ...string) *exec.Cmd {

	var cmd *exec.Cmd

	remoteCmd := []string{"cryosparcm"}

	if cryosparcpath != "" {
		remoteCmd = append(remoteCmd,
			"--dir",
			filepath.Join(cryosparcpath, "cryosparc_master"),
		)
	}

	if commandpath != "" {
		remoteCmd = append(remoteCmd,
			"--cwd",
			commandpath,
		)
	}

	remoteCmd = append(remoteCmd, args...)

	if help || len(args) == 0 {

		cmd = exec.Command(
			filepath.Join(cryosparcpath, "cryosparc_master", "bin", "cryosparcm"),
			"--help",
		)

	} else if len(args) > 1 && args[len(args)-1] == "--help" {

		cmd = exec.Command(
			filepath.Join(cryosparcpath, "cryosparc_master", "bin", "cryosparcm"),
			args...,
		)

	} else {

		cmd = exec.Command(
			"ssh",
			fmt.Sprintf("%s@%s", os.Getenv("USER"), hostname),
			strings.Join(remoteCmd, " "),
		)

	}

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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

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
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

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
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

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
// Create backup
// =============================================================================
func (i *Command) backupCryosparcDatabase(_ func(float64), cryosparcpath, dbpath, backuppath string) error {

	var version string

	date := time.Now().Format("2006-01-02T15:04:05")

	versionData, err := os.ReadFile(
		filepath.Join(cryosparcpath, "cryosparc_master", "version"))

	if err == nil {
		version = strings.TrimSpace(string(versionData))
		version = strings.TrimPrefix(version, "v")
	} else {
		return fmt.Errorf("cryoSPARC version not found: %w", err)
	}

	baseName := filepath.Base(filepath.Clean(dbpath))
	destination := filepath.Join(
		backuppath,
		fmt.Sprintf("%s.v%s.%s.backup", baseName, version, date),
	)

	if err := copyDir(dbpath, destination); err != nil {
		return fmt.Errorf("create backup: %w", err)
	}

	return nil
}

func copyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := in.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(
		dst,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		info.Mode().Perm(),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	_, err = io.Copy(out, in)
	return err
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
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("read gzip archive: %w", err)
	}
	defer func() {
		if err := gz.Close(); err != nil {
			log.Printf("failed to close gzip: %v", err)
		}
	}()

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
	cryosparcpath,
	license,
	hostname,
	dbpath,
	ssdpath,
	arch string,
	baseport uint) *exec.Cmd {

	installScript := filepath.Join(
		cryosparcpath,
		"cryosparc_master",
		"install.sh",
	)

	installArgs := []string{
		installScript,
		"--license", license,
		"--hostname", hostname,
		"--dbpath", dbpath,
		"--ssdpath", ssdpath,
		"--port", strconv.Itoa(int(baseport)),
	}

	// Get local architecture.
	localArch, err := exec.Command("uname", "-m").Output()

	if err != nil {
		localArch = nil
	}

	isLocal := strings.TrimSpace(string(localArch)) == arch
	var args []string

	if isLocal {
		// Run locally.
		args = installArgs
	} else if arch == "aarch64" {
		// Run on ARM nodes.
		args = append([]string{
			"srun",
			"--cluster", "gmerlin7",
			"--partition", "gh-interactive",
			"--time", "0-00:30:00"},
			installArgs...)
	} else {
		// Run on x86 nodes.
		args = append([]string{
			"srun",
			"--cluster", "merlin7",
			"--partition", "interactive",
			"--reservation", "interactive",
			"--time", "0-00:30:00"},
			installArgs...)
	}

	versionData, err := os.ReadFile(
		filepath.Join(cryosparcpath, "cryosparc_master", "version"))

	if err == nil {
		version := strings.TrimSpace(string(versionData))
		version = strings.TrimPrefix(version, "v")

		major, _, ok := strings.Cut(version, ".")
		if !ok {
			major = version
		}

		if majorVersion, err := strconv.Atoi(major); err == nil && majorVersion >= 5 {
			args = append(args, "--ignore-port-conflicts")
		}
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// =============================================================================
// Install Worker
// =============================================================================
func (i *Command) installWorker(cryosparcpath, license, arch string) *exec.Cmd {

	installScript := filepath.Join(
		cryosparcpath,
		"cryosparc_worker",
		"install.sh",
	)

	var cmd *exec.Cmd

	// Get local architecture
	localArch, err := exec.Command("uname", "-m").Output()
	if err != nil {
		localArch = []byte("")
	}

	isLocal := strings.TrimSpace(string(localArch)) == arch

	if isLocal {
		// Run locally
		cmd = exec.Command(
			installScript,
			"--license", license,
		)
	} else if arch == "aarch64" {
		// Run on ARM nodes
		cmd = exec.Command(
			"srun",
			"--cluster", "gmerlin7",
			"--partition", "gh-interactive",
			"--time", "0-00:30:00",
			installScript,
			"--license", license,
		)
	} else {
		// Run on x86 nodes
		cmd = exec.Command(
			"srun",
			"--cluster", "merlin7",
			"--partition", "interactive",
			"--reservation", "interactive",
			"--time", "0-00:30:00",
			installScript,
			"--license", license,
		)
	}

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

// =============================================================================
// Create new lane
// =============================================================================

type ClusterInfo struct {
	Name          string `json:"name"`
	WorkerBinPath string `json:"worker_bin_path"`
	CachePath     string `json:"cache_path"`
	SendCmdTpl    string `json:"send_cmd_tpl"`
	QsubCmdTpl    string `json:"qsub_cmd_tpl"`
	QstatCmdTpl   string `json:"qstat_cmd_tpl"`
	QdelCmdTpl    string `json:"qdel_cmd_tpl"`
	QinfoCmdTpl   string `json:"qinfo_cmd_tpl"`
}

func (i *Command) CreateLane(
	_ func(float64),
	cryosparcpath,
	name,
	cachepath,
	memory,
	time,
	partition,
	gpus,
	cpuspertask,
	cluster string,
) error {

	// Directory where the generated files will be placed.
	laneDir := filepath.Join(cryosparcpath, "lanes", name)

	if err := os.MkdirAll(laneDir, 0755); err != nil {
		return fmt.Errorf("create lane directory: %w", err)
	}

	workerBinPath := filepath.Join(
		cryosparcpath,
		"cryosparc_worker",
		"bin",
		"cryosparcw",
	)

	// ---------------------------------------------------------------------
	// cluster_info.json
	// ---------------------------------------------------------------------

	info := ClusterInfo{
		Name:          name,
		WorkerBinPath: workerBinPath,
		CachePath:     cachepath,

		SendCmdTpl: "{{ command }}",

		QsubCmdTpl: fmt.Sprintf(
			"bash -c 'sbatch --parsable --cluster=%s \"{{ script_path_abs }}\" | cut -d \";\" -f 1'",
			cluster,
		),

		QstatCmdTpl: fmt.Sprintf(
			"squeue --cluster=%s -j {{ cluster_job_id }}",
			cluster,
		),

		QdelCmdTpl: fmt.Sprintf(
			"scancel --cluster=%s {{ cluster_job_id }}",
			cluster,
		),

		QinfoCmdTpl: fmt.Sprintf(
			"sinfo --cluster=%s",
			cluster,
		),
	}

	infoJSON, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal cluster info: %w", err)
	}

	infoPath := filepath.Join(laneDir, "cluster_info.json")

	if err := os.WriteFile(infoPath, infoJSON, 0644); err != nil {
		return fmt.Errorf("write %s: %w", infoPath, err)
	}

	// ---------------------------------------------------------------------
	// cluster_script.sh
	// ---------------------------------------------------------------------

	var script string

	if cluster == "gmerlin7" {
		script = fmt.Sprintf(`#!/usr/bin/env bash

# cryoSPARC cluster submission script template for SLURM
# Lane: %s

{# This template uses jinja2 syntax. #}

# Available variables:
#  script_path_abs={{ script_path_abs }}
#  run_cmd={{ run_cmd }}
#  num_cpu={{ num_cpu }}
#  num_gpu={{ num_gpu }}
#  ram_gb={{ ram_gb }}
#  job_dir_abs={{ job_dir_abs }}
#  project_dir_abs={{ project_dir_abs }}
#  job_log_path_abs={{ job_log_path_abs }}
#  worker_bin_path={{ worker_bin_path }}
#  run_args={{ run_args }}
#  project_uid={{ project_uid }}
#  job_uid={{ job_uid }}
#  job_creator={{ job_creator }}
#  cryosparc_username={{ cryosparc_username }}

#SBATCH --job-name=cryosparc_{{ project_uid }}_{{ job_uid }}
#SBATCH --output={{ job_log_path_abs }}.out
#SBATCH --error={{ job_log_path_abs }}.err
#SBATCH --ntasks=1
#SBATCH --mem={{ (%s * 1000)|int }}M
#SBATCH --time=%s
#SBATCH --partition=%s
#SBATCH --cluster=%s
#SBATCH --gres=gpu:{{ %s }}
#SBATCH --cpus-per-task={{ %s }}

{%%- if num_gpu == 0 %%}
# Use CPU cluster
echo "Error: No GPU requested. Use a CPU lane instead." >&2
exit 1
{%%- else %%}

# Print hostname, for debugging
echo "Job Id: $SLURM_JOBID"
echo "Host:   $SLURM_NODELIST"

module purge

srun {{ run_cmd }}

EXIT_CODE=$?
echo "Exit code: $EXIT_CODE"
exit $?
{%%- endif %%}
`, name, memory, time, partition, cluster, gpus, cpuspertask)
	} else {
		script = fmt.Sprintf(`#!/usr/bin/env bash

# cryoSPARC cluster submission script template for SLURM
# Lane: %s

{# This template uses jinja2 syntax. #}

# Available variables:
#  script_path_abs={{ script_path_abs }}
#  run_cmd={{ run_cmd }}
#  num_cpu={{ num_cpu }}
#  num_gpu={{ num_gpu }}
#  ram_gb={{ ram_gb }}
#  job_dir_abs={{ job_dir_abs }}
#  project_dir_abs={{ project_dir_abs }}
#  job_log_path_abs={{ job_log_path_abs }}
#  worker_bin_path={{ worker_bin_path }}
#  run_args={{ run_args }}
#  project_uid={{ project_uid }}
#  job_uid={{ job_uid }}
#  job_creator={{ job_creator }}
#  cryosparc_username={{ cryosparc_username }}

#SBATCH --job-name=cryosparc_{{ project_uid }}_{{ job_uid }}
#SBATCH --output={{ job_log_path_abs }}.out
#SBATCH --error={{ job_log_path_abs }}.err
#SBATCH --ntasks=1
#SBATCH --mem={{ (%s * 1000)|int }}M
#SBATCH --time=%s
#SBATCH --partition=%s
#SBATCH --cluster=%s
#SBATCH --cpus-per-task={{ %s }}

{%%- if num_gpu > 0 %%}
# Use GPU cluster
echo "Error: GPU requested. Use a GPU lane instead." >&2
exit 1
{%%- else %%}

# Print hostname, for debugging
echo "Job Id: $SLURM_JOBID"
echo "Host:   $SLURM_NODELIST"

module purge

srun {{ run_cmd }}

EXIT_CODE=$?
echo "Exit code: $EXIT_CODE"
exit $?
{%%- endif %%}
`, name, memory, time, partition, cluster, cpuspertask)
	}

	scriptPath := filepath.Join(laneDir, "cluster_script.sh")

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("write %s: %w", scriptPath, err)
	}

	return nil
}

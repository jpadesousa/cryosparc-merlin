package command

import (
	"cryosparc-merlin/ui"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"charm.land/huh/v2"
)

// =============================================================================
// cryosparcm commands Steps
// =============================================================================

func (i *Command) CryosparcmStep(hostname, cryosparcpath string, help bool, args ...string) ui.Step {
	return ui.Step{
		Message:          "",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, help, args...)
		},
	}
}

func (i *Command) CryosparcmStartStep(hostname, cryosparcpath string, help bool) ui.Step {
	return ui.Step{
		Message:          "Starting CryoSPARC instance",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, help, "start")
		},
	}
}

// =============================================================================
// Create CryosparcPath Step
// =============================================================================
func (i *Command) createCryosparcPathStep(cryosparcpath string) ui.Step {
	var action string

	return ui.Step{
		Message:          fmt.Sprintf("Creating CryoSPARC directory %s", cryosparcpath),
		CompletedMessage: fmt.Sprintf("Created CryoSPARC directory %s", cryosparcpath),
		Condition: func() (bool, error) {

			return directoryExistsAndNotEmpty(cryosparcpath)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							cryosparcpath,
						)).
						Options(
							huh.NewOption("Overwrite installation", "overwrite"),
							huh.NewOption("Cancel installation", "cancel"),
						).
						Value(&action),
				),
			)
		},
		Action: func(update func(float64)) error {
			switch action {

			case "cancel":
				return ErrInstallationCancelled

			case "overwrite":
				if err := os.RemoveAll(cryosparcpath); err != nil {
					return err
				}
				return i.createDirectory(update, cryosparcpath)

			default:
				return i.createDirectory(update, cryosparcpath)

			}
		},
	}
}

// =============================================================================
// Check Install Directory
// =============================================================================
func (i *Command) checkInstallDirStep(installDir string) ui.Step {
	var action string

	return ui.Step{
		Message:          fmt.Sprintf("Checking installation directory %s", installDir),
		CompletedMessage: fmt.Sprintf("Created installation directory %s", installDir),
		Condition: func() (bool, error) {

			return directoryExistsAndNotEmpty(installDir)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							installDir,
						)).
						Options(
							huh.NewOption("Overwrite directory", "overwrite"),
							huh.NewOption("Cancel installation", "cancel"),
						).
						Value(&action),
				),
			)
		},
		Action: func(update func(float64)) error {
			switch action {

			case "cancel":
				return ErrInstallationCancelled

			case "overwrite":
				if err := os.RemoveAll(installDir); err != nil {
					return err
				}

				return i.createDirectory(update, installDir)

			default:
				return i.createDirectory(update, installDir)

			}
		},
	}
}

// =============================================================================
// Create DbPath Step
// =============================================================================
func (i *Command) createDbPathStep(dbpath string) ui.Step {
	var action string

	return ui.Step{
		Message:          fmt.Sprintf("Creating database directory %s", dbpath),
		CompletedMessage: fmt.Sprintf("Created database directory %s", dbpath),
		Condition: func() (bool, error) {

			return directoryExistsAndNotEmpty(dbpath)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							dbpath,
						)).
						Options(
							huh.NewOption("Overwrite database", "overwrite"),
							huh.NewOption(
								"Continue with existing database\n(warning: database cannot be shared between instances)", "continue"),
							huh.NewOption("Cancel installation", "cancel"),
						).
						Value(&action),
				),
			)
		},
		Skip: func() bool {
			return action == "continue"
		},
		SkipMessage: "Using existing database",
		Action: func(update func(float64)) error {
			switch action {
			case "cancel":
				return ErrInstallationCancelled

			case "overwrite":
				if err := os.RemoveAll(dbpath); err != nil {
					return err
				}
				return i.createDirectory(update, dbpath)

			case "continue":
				// Should never be reached because Skip() handled it.
				return nil

			default:
				return i.createDirectory(update, dbpath)

			}
		},
	}
}

// =============================================================================
// Download CryoSPARC Step
// =============================================================================
var newDownloadDir string

func (i *Command) downloadCryosparcStep(downloaddir, release, license, installation, arch string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Downloading CryoSPARC %s v%s arch=%s", installation, release, arch),
		CompletedMessage: fmt.Sprintf("Downloaded CryoSPARC %s v%s arch=%s", installation, release, arch),
		Skip: func() bool {

			basePath := fmt.Sprintf(
				"/data/project/cls/shared/software/cryosparc/installations/v%s/%s",
				release,
				arch,
			)

			path := filepath.Join(basePath, fmt.Sprintf(
				"cryosparc_%s.tar.gz",
				installation,
			))

			_, err := os.Stat(path)

			if err == nil {
				newDownloadDir = basePath
				return true
			}
			return false
		},
		SkipMessage:     fmt.Sprintf("CryoSPARC %s v%s arch=%s was already downloaded", installation, release, arch),
		ShowProgressBar: true,
		Action: func(update func(float64)) error {

			return i.downloadCryosparc(update, downloaddir, release, license, installation, arch)
		},
	}
}

// =============================================================================
// Extract Archive Step
// =============================================================================
func (i *Command) ExtractArchiveStep(extractdir, release, installation, arch string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Extracting CryoSPARC %s v%s arch=%s", installation, release, arch),
		CompletedMessage: fmt.Sprintf("Extracted CryoSPARC %s v%s arch=%s", installation, release, arch),
		Action: func(update func(float64)) error {

			var downloaddir string

			if newDownloadDir != "" {
				downloaddir = newDownloadDir
			} else {
				downloaddir = extractdir
			}

			return i.extractArchive(update, downloaddir, extractdir, fmt.Sprintf("cryosparc_%s.tar.gz", installation))
		},
	}
}

// =============================================================================
// Install Master Step
// =============================================================================
func (i *Command) InstallMasterStep(cryosparcpath, release, license, hostname, dbpath, ssdpath, arch string, baseport uint) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Installing CryoSPARC master v%s arch=%s", release, arch),
		CompletedMessage: fmt.Sprintf("Installed CryoSPARC master v%s arch=%s", release, arch),
		Exec: func() *exec.Cmd {

			return i.installMaster(cryosparcpath, license, hostname, dbpath, ssdpath, arch, baseport)
		},
	}
}

// =============================================================================
// Install Worker Step
// =============================================================================
func (i *Command) InstallWorkerStep(cryosparcpath, release, license, arch string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Installing CryoSPARC worker v%s arch=%s", release, arch),
		CompletedMessage: fmt.Sprintf("Installed CryoSPARC worker v%s arch=%s", release, arch),
		Exec: func() *exec.Cmd {

			return i.installWorker(cryosparcpath, license, arch)
		},
	}
}

// =============================================================================
// Replace license ID in config.sh Step
// =============================================================================
func (i *Command) replaceLicenseIDStep(installDir, license, installation string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Checking CryoSPARC %s License ID", installation),
		CompletedMessage: fmt.Sprintf("Checked CryoSPARC %s License ID", installation),
		Action: func(update func(float64)) error {

			return i.replaceLicenseID(update, installDir, license)
		},
	}
}

// =============================================================================
// Lanes Steps
// =============================================================================
func (i *Command) RunLanesCreateStep(
	cryosparcpath,
	name,
	cachepath,
	memory,
	time,
	partition,
	gpus,
	cpuspertask,
	cluster string) ui.Step {
	return ui.Step{
		Message: fmt.Sprintf("Creating CryoSPARC lane: '%s' (%s)",
			name,
			filepath.Join(cryosparcpath, "lanes"),
		),
		CompletedMessage: fmt.Sprintf("Created CryoSPARC lane: '%s' (%s)",
			name,
			filepath.Join(cryosparcpath, "lanes"),
		),
		Action: func(update func(float64)) error {
			return i.CreateLane(
				update,
				cryosparcpath,
				name,
				cachepath,
				memory,
				time,
				partition,
				gpus,
				cpuspertask,
				cluster)
		},
	}
}

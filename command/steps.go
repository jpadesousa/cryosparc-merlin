package command

import (
	"cryosparc-merlin/ui"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"charm.land/huh/v2"
)

// =============================================================================
// cryosparcm commands Steps
// =============================================================================
func (i *Command) CryosparcmStatusStep(remotehost, homedir string) ui.Step {
	return ui.Step{
		Message:          "Checking CryoSPARC instance status",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(remotehost, homedir, "status")
		},
	}
}

func (i *Command) CryosparcmStartStep(remotehost, homedir string) ui.Step {
	return ui.Step{
		Message:          "Starting CryoSPARC instance",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(remotehost, homedir, "start")
		},
	}
}

func (i *Command) CryosparcmStopStep(remotehost, homedir string) ui.Step {
	return ui.Step{
		Message:          "Stopping CryoSPARC instance",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(remotehost, homedir, "stop")
		},
	}
}

func (i *Command) CryosparcmRestartStep(remotehost, homedir string) ui.Step {
	return ui.Step{
		Message:          "Restarting CryoSPARC instance",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(remotehost, homedir, "restart")
		},
	}
}

func (i *Command) CryosparcmCreateUserStep(remotehost, homedir, email, username, firstName, lastName string) ui.Step {

	var password string

	return ui.Step{
		Message:          "Creating CryoSPARC user",
		CompletedMessage: "",
		Condition: func() (bool, error) {

			return true, nil
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Email").
						Validate(func(s string) error {
							if s == "" {
								return errors.New("Email is required")
							}
							return nil
						}).
						Value(&email),

					huh.NewInput().
						Title("Password").
						EchoMode(huh.EchoModePassword).
						Validate(func(s string) error {
							if s == "" {
								return errors.New("Password is required")
							}
							return nil
						}).
						Value(&password),

					huh.NewInput().
						Title("Username").
						Value(&username),

					huh.NewInput().
						Title("First Name").
						Value(&firstName),

					huh.NewInput().
						Title("Last Name").
						Value(&lastName),
				),
			)
		},
		Exec: func() *exec.Cmd {

			return cryosparcmCmd(
				remotehost,
				homedir,
				"createuser",
				"--email", email,
				"--password", password,
				"--username", username,
				"--firstname", firstName,
				"--lastname", lastName,
			)
		},
	}
}

func (i *Command) CryosparcmResetPasswordStep(remotehost, homedir, email string) ui.Step {

	var password string

	return ui.Step{
		Message:          "Resetting CryoSPARC user password",
		CompletedMessage: "",
		Condition: func() (bool, error) {
			return true, nil
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Email").
						Validate(func(s string) error {
							if s == "" {
								return errors.New("Email is required")
							}
							return nil
						}).
						Value(&email),

					huh.NewInput().
						Title("Password").
						EchoMode(huh.EchoModePassword).
						Validate(func(s string) error {
							if s == "" {
								return errors.New("Password is required")
							}
							return nil
						}).
						Value(&password),
				),
			)
		},
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(
				remotehost,
				homedir,
				"resetpassword",
				"--email", email,
				"--password", password,
			)
		},
	}
}

// func (i *Command) CryosparcmUpdateStep(remotehost, homedir string, args ...string) ui.Step {
// 	return ui.Step{
// 		Message:          "Running CryoSPARC update command",
// 		CompletedMessage: "",
// 		Exec: func() *exec.Cmd {
// 			cmdArgs := append([]string{"update"}, args...)
// 			return cryosparcmCmd(remotehost, homedir, cmdArgs...)
// 		},
// 	}
// }

func (i *Command) CryosparcmUpdateStep(
	remotehost, homedir, version string,
	check, list, override, downloadOnly, skipDownload bool,
) ui.Step {
	return ui.Step{
		Message:          "",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			args := []string{"update"}

			if version != "" {
				args = append(args, "--version="+version)
			}
			if check {
				args = append(args, "--check")
			}
			if list {
				args = append(args, "--list")
			}
			if override {
				args = append(args, "--override")
			}
			if downloadOnly {
				args = append(args, "--download-only")
			}
			if skipDownload {
				args = append(args, "--skip-download")
			}

			return cryosparcmCmd(remotehost, homedir, args...)
		},
	}
}

func (i *Command) CryosparcmPatchStep(
	remotehost, homedir string,
	install, download, check, yes, force bool,
) ui.Step {
	return ui.Step{
		Message:          "",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			args := []string{"patch"}

			if install {
				args = append(args, "--install")
			}
			if download {
				args = append(args, "--download")
			}
			if check {
				args = append(args, "--check")
			}
			if yes {
				args = append(args, "--yes")
			}
			if force {
				args = append(args, "--force")
			}

			return cryosparcmCmd(remotehost, homedir, args...)
		},
	}
}

// =============================================================================
// Create HomeDir Step
// =============================================================================
func (i *Command) createHomeDirStep(homedir string) ui.Step {
	var action string

	return ui.Step{
		Message:          fmt.Sprintf("Creating home directory %s", homedir),
		CompletedMessage: fmt.Sprintf("Created home directory %s", homedir),
		Condition: func() (bool, error) {

			return directoryExistsAndNotEmpty(homedir)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							homedir,
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
				if err := os.RemoveAll(homedir); err != nil {
					return err
				}
				return i.createDirectory(update, homedir)

			default:
				return i.createDirectory(update, homedir)

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
// Create DbDir Step
// =============================================================================
func (i *Command) createDbDirStep(dbdir string) ui.Step {
	var action string

	return ui.Step{
		Message:          fmt.Sprintf("Creating database directory %s", dbdir),
		CompletedMessage: fmt.Sprintf("Created database directory %s", dbdir),
		Condition: func() (bool, error) {

			return directoryExistsAndNotEmpty(dbdir)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							dbdir,
						)).
						Options(
							huh.NewOption("Overwrite database", "overwrite"),
							huh.NewOption("Continue with existing database", "continue"),
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
				if err := os.RemoveAll(dbdir); err != nil {
					return err
				}
				return i.createDirectory(update, dbdir)

			case "continue":
				// Should never be reached because Skip() handled it.
				return nil

			default:
				return i.createDirectory(update, dbdir)

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
func (i *Command) InstallMasterStep(homedir, release, license, remotehost, dbdir, arch string, baseport uint) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Installing CryoSPARC master v%s arch=%s", release, arch),
		CompletedMessage: fmt.Sprintf("Installed CryoSPARC master v%s arch=%s", release, arch),
		Exec: func() *exec.Cmd {

			return i.installMaster(homedir, license, remotehost, dbdir, baseport)
		},
	}
}

// =============================================================================
// Install Worker Step
// =============================================================================
func (i *Command) InstallWorkerStep(homedir, release, license, arch string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Installing CryoSPARC worker v%s arch=%s", release, arch),
		CompletedMessage: fmt.Sprintf("Installed CryoSPARC worker v%s arch=%s", release, arch),
		Exec: func() *exec.Cmd {

			return i.installWorker(homedir, license)
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

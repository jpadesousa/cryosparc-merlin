package command

import (
	"cryosparc-merlin/ui"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"charm.land/huh/v2"
)

// =============================================================================
// cryosparcm commands Steps
// =============================================================================

func (i *Command) CryosparcmStep(hostname, cryosparcpath string, help bool, args ...string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Running cryosparcm %s", strings.Join(args, " ")),
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, "", help, args...)
		},
	}
}

func (i *Command) CryosparcmStartStep(hostname, cryosparcpath string, help bool) ui.Step {
	return ui.Step{
		Message:          "Starting CryoSPARC instance",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, "", help, "start")
		},
	}
}

func (i *Command) CryosparcmCreateUserStep(hostname, cryosparcpath string, help bool) ui.Step {

	var password, email, username, firstName, lastName string

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
								return errors.New("email is required")
							}
							return nil
						}).
						Value(&email),

					huh.NewInput().
						Title("Password").
						EchoMode(huh.EchoModePassword).
						Validate(func(s string) error {
							if s == "" {
								return errors.New("password is required")
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

			var args []string

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
					args = append(args, "user", "create")
				}
			} else {
				args = append(args, "createuser")
			}

			args = append(args,
				"--email", email,
				"--password", password,
				"--username", username,
				"--firstname", firstName,
				"--lastname", lastName)

			return cryosparcmCmd(
				hostname,
				cryosparcpath,
				"",
				help,
				args...)
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
				return ErrCancelled

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
				return ErrCancelled

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
// Backup CryoSPARC database
// =============================================================================

func (i *Command) backupCryosparcDatabaseStep(cryosparcpath, dbpath, backupdir string) ui.Step {
	return ui.Step{
		Message:          fmt.Sprintf("Backing up CryoSPARC database (%s)", dbpath),
		CompletedMessage: fmt.Sprintf("Backed up CryoSPARC database (%s)", backupdir),
		Action: func(update func(float64)) error {
			return i.backupCryosparcDatabase(update, cryosparcpath, dbpath, backupdir)
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
				return ErrCancelled

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
		Message: fmt.Sprintf("Installing CryoSPARC master v%s arch=%s", release, arch),
		// CompletedMessage: fmt.Sprintf("Installed CryoSPARC master v%s arch=%s", release, arch),
		CompletedMessage: "",
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
		Message: fmt.Sprintf("Installing CryoSPARC worker v%s arch=%s", release, arch),
		// CompletedMessage: fmt.Sprintf("Installed CryoSPARC worker v%s arch=%s", release, arch),
		CompletedMessage: "",
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
			filepath.Join(cryosparcpath, "lanes", name),
		),
		CompletedMessage: fmt.Sprintf("Created CryoSPARC lane: '%s' (%s)",
			name,
			filepath.Join(cryosparcpath, "lanes", name),
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

func (i *Command) RunLanesInstallStep(
	hostname,
	cryosparcpath,
	lanepath,
	info,
	script string,
	help bool) ui.Step {
	return ui.Step{
		Message:          "Installing CryoSPARC lane",
		CompletedMessage: "",
		Exec: func() *exec.Cmd {

			cmdArgs := append([]string{"cluster", "connect"},
				"--info", info,
				"--script", script,
			)

			return cryosparcmCmd(
				hostname,
				cryosparcpath,
				lanepath,
				help,
				cmdArgs...,
			)

		},
	}
}

func (i *Command) RunLanesRemoveStep(hostname, cryosparcpath, name string, help bool) ui.Step {
	return ui.Step{
		Message: fmt.Sprintf("Removing CryoSPARC lane (%s)",
			name,
		),
		CompletedMessage: "",
		Exec: func() *exec.Cmd {

			cmdArgs := append([]string{"cluster", "remove"}, name)

			return cryosparcmCmd(hostname, cryosparcpath, "", help, cmdArgs...)
		},
	}
}

func (i *Command) RunLanesCreateDefaultSteps(cryosparcpath, cachepath, memory, gpus, cpuspertask string) []ui.Step {

	var arch string

	archData, err := os.ReadFile(
		filepath.Join(cryosparcpath, "cryosparc_master", "arch"))

	if err == nil {
		arch = strings.TrimSpace(string(archData))
	} else {
		arch = ""
	}

	steps := map[string]ui.Step{
		"cpu-hourly": i.RunLanesCreateStep(
			cryosparcpath,
			"cpu-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"hourly",
			gpus,
			cpuspertask,
			"merlin7",
		),
		"cpu-daily": i.RunLanesCreateStep(
			cryosparcpath,
			"cpu-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"daily",
			gpus,
			cpuspertask,
			"merlin7",
		),
	}

	if arch == "aarch64" {

		steps["gh-hourly"] = i.RunLanesCreateStep(
			cryosparcpath,
			"gh-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"gh-hourly",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["gh-daily"] = i.RunLanesCreateStep(
			cryosparcpath,
			"gh-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"gh-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["gh-daily-4h"] = i.RunLanesCreateStep(
			cryosparcpath,
			"gh-daily-4h",
			cachepath,
			memory,
			"00-04:00:00",
			"gh-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["gh-general-2d"] = i.RunLanesCreateStep(
			cryosparcpath,
			"gh-general-2d",
			cachepath,
			memory,
			"02-00:00:00",
			"gh-general",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["gh-general-2d-100Gb"] = i.RunLanesCreateStep(
			cryosparcpath,
			"gh-general-2d-100Gb",
			cachepath,
			"100",
			"02-00:00:00",
			"gh-general",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

	} else {
		steps["a100-hourly"] = i.RunLanesCreateStep(
			cryosparcpath,
			"a100-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"a100-hourly",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["a100-daily"] = i.RunLanesCreateStep(
			cryosparcpath,
			"a100-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"a100-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["a100-daily-4h"] = i.RunLanesCreateStep(
			cryosparcpath,
			"a100-daily-4h",
			cachepath,
			memory,
			"00-04:00:00",
			"a100-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["a100-general-2d"] = i.RunLanesCreateStep(
			cryosparcpath,
			"a100-general-2d",
			cachepath,
			memory,
			"02-00:00:00",
			"a100-general",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

		steps["a100-general-2d-100Gb"] = i.RunLanesCreateStep(
			cryosparcpath,
			"a100-general-2d-100Gb",
			cachepath,
			"100",
			"02-00:00:00",
			"a100-general",
			gpus,
			cpuspertask,
			"gmerlin7",
		)

	}

	keys := make([]string, 0, len(steps))
	for key := range steps {
		keys = append(keys, key)
	}

	stepsList := make([]ui.Step, 0, len(steps)+1)
	stepsList = append(
		stepsList,
		i.RunLanesCreateDefaultConfirmStep(cryosparcpath, keys),
	)

	for key := range steps {
		stepsList = append(stepsList, steps[key])
	}

	return stepsList
}

func (i *Command) RunLanesCreateDefaultConfirmStep(
	cryosparcPath string,
	stepNames []string,
) ui.Step {
	var proceed bool

	return ui.Step{
		Message:          "Creating all default lanes",
		CompletedMessage: "",
		Condition: func() (bool, error) {
			return true, nil
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf(
							"This will create all default lanes in %s:\n\n%s\n\nDo you want to proceed?",
							cryosparcPath,
							strings.Join(stepNames, "\n"),
						)).
						Affirmative("Confirm").
						Negative("Cancel").
						Value(&proceed),
				),
			)
		},
		Action: func(update func(float64)) error {
			if !proceed {
				return ErrCancelled
			}

			return nil
		},
	}
}

func (i *Command) RunLanesInstallSteps(hostname, cryosparcpath, info, script, name string, all, help bool) []ui.Step {

	var lanes []string

	steps := make([]ui.Step, 0)

	if info == "" && script == "" {

		if name != "" {

			lanes = append(lanes, filepath.Join(cryosparcpath, "lanes", name))

		} else if all {

			var err error
			lanes, err = filepath.Glob(
				filepath.Join(cryosparcpath, "lanes", "*"),
			)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return nil
			}

		} else {
			fmt.Println("no lanes selected. Use '--name' to select a lane to install or '--all' to install all lanes stored in the CryoSPARC directory")
			return nil
		}

		for _, lanepath := range lanes {

			info = filepath.Join(lanepath, "cluster_info.json")
			script = filepath.Join(lanepath, "cluster_script.sh")

			steps = append(steps,
				i.RunLanesInstallStep(
					hostname,
					cryosparcpath,
					lanepath,
					info,
					script,
					help,
				),
			)
		}

	} else {
		steps = append(steps,
			i.RunLanesInstallStep(
				hostname,
				cryosparcpath,
				"",
				info,
				script,
				help,
			),
		)
	}

	return steps
}

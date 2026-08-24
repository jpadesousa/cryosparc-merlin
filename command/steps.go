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
func (i *Command) CryosparcmStep(
	hostname,
	cryosparcpath string,
	help bool,
	args ...string) ui.Step {

	return ui.Step{
		Message: strings.TrimSpace(
			fmt.Sprintf("Running cryosparcm %s", strings.Join(args, " "))) +
			fmt.Sprintf(" (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, "", help, args...)
		},
	}
}

func (i *Command) CryosparcmStartStep(
	hostname,
	cryosparcpath string,
	help bool) ui.Step {

	return ui.Step{
		Message:          fmt.Sprintf("Starting CryoSPARC (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, "", help, "start")
		},
	}
}

func (i *Command) CryosparcmStopStep(
	hostname,
	cryosparcpath string,
	help bool) ui.Step {

	return ui.Step{
		Message:          fmt.Sprintf("Stopping CryoSPARC (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {
			return cryosparcmCmd(hostname, cryosparcpath, "", help, "stop")
		},
	}
}

func (i *Command) CryosparcmCreateUserStep(
	hostname,
	cryosparcpath,
	email,
	username,
	firstname,
	lastname string,
	help bool) ui.Step {

	var password string

	return ui.Step{
		Message:          fmt.Sprintf("Creating CryoSPARC user (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Condition: func() (bool, error) {
			return true, nil
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewNote().
						Title("Create CryoSPARC user"),

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
						Title("Password (different from your PSI password)").
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
						Value(&firstname),

					huh.NewInput().
						Title("Last Name").
						Value(&lastname),
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

				if majorVersion, err := strconv.Atoi(major); err == nil &&
					majorVersion >= 5 {
					args = append(args, "user", "create")
				} else {
					args = append(args, "createuser")
				}
			} else {
				args = append(args, "createuser")
			}

			args = append(args,
				"--email", email,
				"--password", password,
				"--username", username,
				"--firstname", firstname,
				"--lastname", lastname)

			return cryosparcmCmd(
				hostname,
				cryosparcpath,
				"",
				help,
				args...)
		},
	}
}

func (i *Command) CryosparcmChangePortStep(
	hostname,
	cryosparcpath string,
	port, start, end, count uint,
	help bool) ui.Step {

	return ui.Step{
		Message:          fmt.Sprintf("Changing base port (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {

			if port == 0 {
				i.setCryosparcBasePort(hostname, start, end, count)
			}

			return cryosparcmCmd(
				hostname,
				cryosparcpath,
				"",
				help,
				"changeport",
				strconv.FormatUint(uint64(i.cfg.BasePort), 10),
			)
		},
	}
}

// =============================================================================
// cryosparcm commands Steps
// =============================================================================
func (i *Command) CryosparcwStep(
	cryosparcpath string,
	arch string,
	help bool,
	args ...string) ui.Step {

	return ui.Step{
		Message: strings.TrimSpace(
			fmt.Sprintf("Running cryosparcw %s", strings.Join(args, " "))) +
			fmt.Sprintf(" (instance: %s)", cryosparcpath),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {
			return cryosparcwCmd(cryosparcpath, arch, help, args...)
		},
	}
}

// =============================================================================
// Check Install Directory
// =============================================================================
func (i *Command) checkInstallDirStep(oldDir, newDir string) ui.Step {
	var action string

	return ui.Step{
		Message: fmt.Sprintf(
			"Preparing installation directory %s", oldDir),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Created installation directory %s", newDir)
		},
		Condition: func() (bool, error) {
			return directoryExistsAndNotEmpty(oldDir)
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							oldDir,
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

				if err := os.RemoveAll(oldDir); err != nil {
					return err
				}

				return i.createDirectory(update, newDir)

			default:
				return i.createDirectory(update, newDir)

			}
		},
	}
}

// =============================================================================
// Backup CryoSPARC database
// =============================================================================
func (i *Command) backupCryosparcDatabaseStep(
	cryosparcpath,
	dbpath,
	backupdir string) ui.Step {
	return ui.Step{
		Message: fmt.Sprintf(
			"Backing up CryoSPARC database (%s)", dbpath),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Backed up CryoSPARC database (destination: %s)", backupdir)
		},
		Skip: func() bool {
			exists, _ := directoryExistsAndNotEmpty(dbpath)
			return !exists
		},
		SkipMessage: fmt.Sprintf(
			"Database backup skipped: Previous database is empty or does not exist (%s)",
			dbpath),
		Action: func(update func(float64)) error {
			return i.backupCryosparcDatabase(
				update,
				cryosparcpath,
				dbpath,
				backupdir)
		},
	}
}

// =============================================================================
// Create DbPath Step
// =============================================================================
func (i *Command) createDbPathStep(dbpath string) ui.Step {
	var action string

	return ui.Step{
		Message: fmt.Sprintf(
			"Creating database directory %s", dbpath),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Created database directory %s", dbpath)
		},
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
				if err := os.RemoveAll(dbpath); err != nil {
					return err
				}
				return i.createDirectory(update, dbpath)

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

func (i *Command) downloadCryosparcStep(
	downloaddir,
	release,
	license,
	installation,
	arch string) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf(
			"Downloading CryoSPARC %s v%s arch=%s", installation, release, arch),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Downloaded CryoSPARC %s v%s arch=%s", installation, release, arch)
		},
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
		SkipMessage: fmt.Sprintf(
			"CryoSPARC %s v%s arch=%s was already downloaded",
			installation,
			release,
			arch),
		ShowProgressBar: true,
		Action: func(update func(float64)) error {

			return i.downloadCryosparc(
				update,
				downloaddir,
				release,
				license,
				installation,
				arch)
		},
	}
}

// =============================================================================
// Extract Archive Step
// =============================================================================
func (i *Command) ExtractArchiveStep(
	extractdir,
	release,
	installation,
	arch string) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf(
			"Extracting CryoSPARC %s v%s arch=%s", installation, release, arch),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Extracted CryoSPARC %s v%s arch=%s", installation, release, arch)
		},
		Action: func(update func(float64)) error {

			var downloaddir string

			if newDownloadDir != "" {
				downloaddir = newDownloadDir
			} else {
				downloaddir = extractdir
			}

			return i.extractArchive(
				update,
				downloaddir,
				extractdir,
				fmt.Sprintf("cryosparc_%s.tar.gz", installation))
		},
	}
}

func (i *Command) RenameWorkerDirectoryStep(
	cryosparcpath,
	arch string) ui.Step {

	return ui.Step{
		Message: "Renaming worker directory",
		CompletedMessage: func() string {
			return fmt.Sprintf("Directory renamed to: %s",
				filepath.Join(cryosparcpath,
					fmt.Sprintf("cryosparc_worker_%s", arch)))
		},
		Action: func(update func(float64)) error {
			return i.renameDirectory(
				filepath.Join(cryosparcpath, "cryosparc_worker"),
				filepath.Join(cryosparcpath,
					fmt.Sprintf("cryosparc_worker_%s", arch)),
			)
		},
	}
}

// =============================================================================
// Install Master Step
// =============================================================================
func (i *Command) InstallMasterStep(
	cryosparcpath,
	release,
	license,
	hostname,
	dbpath,
	ssdpath,
	arch string) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf(
			"Installing CryoSPARC master v%s arch=%s", release, arch),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {

			return i.installMaster(
				cryosparcpath,
				license,
				hostname,
				dbpath,
				ssdpath,
				arch,
				i.cfg.BasePort)
		},
	}
}

// =============================================================================
// Install Worker Step
// =============================================================================
func (i *Command) InstallWorkerStep(
	cryosparcpath,
	release,
	license,
	arch string) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf(
			"Installing CryoSPARC worker v%s arch=%s", release, arch),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {

			return i.installWorker(
				cryosparcpath,
				license,
				arch)
		},
	}

}

// =============================================================================
// Replace license ID in config.sh Step
// =============================================================================
func (i *Command) replaceLicenseIDStep(
	installDir,
	license,
	installation string) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf(
			"Checking CryoSPARC %s License ID", installation),
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Checked CryoSPARC %s License ID", installation)
		},
		Action: func(update func(float64)) error {
			return i.replaceLicenseID(update, installDir, license)
		},
	}
}

// =============================================================================
// Lanes Steps (single lanes)
// =============================================================================

// create
func (i *Command) LanesCreateStep(
	cryosparcpath,
	name,
	cachepath,
	memory,
	time,
	partition,
	gpus,
	cpuspertask,
	cluster,
	arch string) ui.Step {

	var action string

	return ui.Step{
		Message: fmt.Sprintf("Creating CryoSPARC lane: '%s' (%s)",
			name,
			filepath.Join(cryosparcpath, "lanes", name),
		),
		CompletedMessage: func() string {
			return fmt.Sprintf("Created CryoSPARC lane: '%s' (%s)",
				name,
				filepath.Join(cryosparcpath, "lanes", name),
			)
		},
		Condition: func() (bool, error) {
			return directoryExistsAndNotEmpty(
				filepath.Join(cryosparcpath, "lanes", name))
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title(fmt.Sprintf(
							"%q already exists",
							filepath.Join(cryosparcpath, "lanes", name),
						)).
						Options(
							huh.NewOption("Overwrite lane", "overwrite"),
							huh.NewOption("Skip lane", "skip"),
						).
						Value(&action),
				),
			)
		},
		Skip: func() bool {
			return action == "skip"
		},
		SkipMessage: fmt.Sprintf(
			"Skipped creating lane %s",
			filepath.Join(cryosparcpath, "lanes", name)),
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
				cluster,
				arch)
		},
	}
}

// install
func (i *Command) LanesInstallStep(
	hostname,
	cryosparcpath,
	lanepath,
	info,
	script string,
	help bool) ui.Step {
	return ui.Step{
		Message: fmt.Sprintf("Installing CryoSPARC lane (instance: %s)",
			cryosparcpath,
		),
		CompletedMessage: nil,
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

// remove
func (i *Command) LanesRemoveStep(
	hostname,
	cryosparcpath,
	name string,
	help bool) ui.Step {

	return ui.Step{
		Message: fmt.Sprintf("Removing CryoSPARC lane from database (%s)",
			name,
		),
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {

			cmdArgs := append([]string{"cluster", "remove"}, name)

			return cryosparcmCmd(hostname, cryosparcpath, "", help, cmdArgs...)
		},
	}
}

// =============================================================================
// Lanes Steps (multiple lanes)
// =============================================================================

// create default

// LanesCreateDefaultSteps creates the default lanes created to run CryoSPARC
// in the Merlin cluster
func (i *Command) LanesCreateDefaultSteps(
	cryosparcpath,
	cachepath,
	memory,
	gpus,
	cpuspertask,
	arch string) []ui.Step {

	steps := make(map[string]ui.Step)

	if arch == "aarch64" {

		steps["gh-hourly"] = i.LanesCreateStep(
			cryosparcpath,
			"gh-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"gh-hourly",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["gh-daily"] = i.LanesCreateStep(
			cryosparcpath,
			"gh-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"gh-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["gh-daily-4h"] = i.LanesCreateStep(
			cryosparcpath,
			"gh-daily-4h",
			cachepath,
			memory,
			"00-04:00:00",
			"gh-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["gh-general-2d"] = i.LanesCreateStep(
			cryosparcpath,
			"gh-general-2d",
			cachepath,
			memory,
			"02-00:00:00",
			"gh-general",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["gh-general-2d-100Gb"] = i.LanesCreateStep(
			cryosparcpath,
			"gh-general-2d-100Gb",
			cachepath,
			"100",
			"02-00:00:00",
			"gh-general",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

	} else {
		steps["cpu-hourly"] = i.LanesCreateStep(
			cryosparcpath,
			"cpu-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"hourly",
			gpus,
			cpuspertask,
			"merlin7",
			arch,
		)

		steps["cpu-daily"] = i.LanesCreateStep(
			cryosparcpath,
			"cpu-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"daily",
			gpus,
			cpuspertask,
			"merlin7",
			arch,
		)

		steps["a100-hourly"] = i.LanesCreateStep(
			cryosparcpath,
			"a100-hourly",
			cachepath,
			memory,
			"00-01:00:00",
			"a100-hourly",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["a100-daily"] = i.LanesCreateStep(
			cryosparcpath,
			"a100-daily",
			cachepath,
			memory,
			"01-00:00:00",
			"a100-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["a100-daily-4h"] = i.LanesCreateStep(
			cryosparcpath,
			"a100-daily-4h",
			cachepath,
			memory,
			"00-04:00:00",
			"a100-daily",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["a100-general-2d"] = i.LanesCreateStep(
			cryosparcpath,
			"a100-general-2d",
			cachepath,
			memory,
			"02-00:00:00",
			"a100-general",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

		steps["a100-general-2d-100Gb"] = i.LanesCreateStep(
			cryosparcpath,
			"a100-general-2d-100Gb",
			cachepath,
			"100",
			"02-00:00:00",
			"a100-general",
			gpus,
			cpuspertask,
			"gmerlin7",
			arch,
		)

	}

	keys := make([]string, 0, len(steps))
	for key := range steps {
		keys = append(keys, key)
	}

	stepsList := make([]ui.Step, 0, len(steps)+1)
	stepsList = append(
		stepsList,
		i.LanesCreateDefaultConfirmStep(cryosparcpath, keys),
	)

	for key := range steps {
		stepsList = append(stepsList, steps[key])
	}

	return stepsList
}

// create default (confirm step)

// LanesCreateDefaultConfirmStep is used to confirm if the creation of all
// default lanes into cryosparc path can proceed.
func (i *Command) LanesCreateDefaultConfirmStep(
	cryosparcPath string,
	stepNames []string,
) ui.Step {
	var proceed bool

	return ui.Step{
		Message:          "Creating all default lanes",
		CompletedMessage: nil,
		Condition: func() (bool, error) {
			return true, nil
		},
		Prompt: func() *huh.Form {
			return huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf(
							"This will create all default lanes "+
								"in %s:\n\n%s\n\nDo you want to proceed?",
							filepath.Join(cryosparcPath, "lanes"),
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

// install

// LanesInstallSteps combines the installation of multiple lanes using the
// previously defined single lane install (RunLanesInstallStep).
func (i *Command) LanesInstallSteps(
	hostname,
	cryosparcpath,
	info,
	script,
	name string,
	all,
	help bool) []ui.Step {

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
			fmt.Println("no lanes selected. Use '--name <lane name>' " +
				"to select a lane to install or '--all' to install all " +
				"lanes stored in the CryoSPARC directory.")
			return nil
		}

		for _, lanepath := range lanes {

			info = filepath.Join(lanepath, "cluster_info.json")
			script = filepath.Join(lanepath, "cluster_script.sh")

			steps = append(steps,
				i.LanesInstallStep(
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
			i.LanesInstallStep(
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

// =============================================================================
// HTTP link
// =============================================================================
func (i *Command) HttpUrlStep(hostname string) ui.Step {

	return ui.Step{
		Message: "",
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"CryoSPARC instance can be accessed on: http://%s:%d",
				hostname,
				i.cfg.BasePort,
			)
		},
		Action: func(update func(float64)) error {
			return nil
		},
	}

}

// =============================================================================
// Set CryoSPARC base port
// =============================================================================
func (i *Command) setCryosparcBasePortStep(
	hostname string,
	port, start, end, count uint) ui.Step {

	return ui.Step{
		Message: "Setting CryoSPARC base port",
		CompletedMessage: func() string {
			return fmt.Sprintf(
				"Set CryoSPARC base port to %d",
				i.cfg.BasePort,
			)
		},
		Skip: func() bool {
			return port != 0
		},
		SkipMessage: fmt.Sprintf(
			"CryoSPARC base port was already provided and set to %d", port),
		Action: func(update func(float64)) error {

			return i.setCryosparcBasePort(hostname, start, end, count)

		},
	}
}

// =============================================================================
// Check running instances
// =============================================================================
func (i *Command) CheckInstancesStep(hostname string) ui.Step {

	return ui.Step{
		Message:          "Checking running CryoSPARC instances",
		CompletedMessage: nil,
		Exec: func() *exec.Cmd {

			return i.checkInstances(hostname)

		},
	}
}

package command

import (
	"cryosparc-merlin/ui"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	portStart uint = 61000 // Start port range to assign new base port
	portEnd   uint = 62000 // End port range to assign new base port
	portCount uint = 10    // Number of contiguous ports available
)

// install complete
func (i *Command) RunInstallComplete() error {

	// ====================================================================== //
	// Only after version 5 that CryoSPARC allows to have different cryosparc
	// worker installations
	// ====================================================================== //
	version := strings.TrimSpace(i.cfg.Version)

	major, _, ok := strings.Cut(version, ".")
	if !ok {
		major = version
	}

	majorVersion, err := strconv.Atoi(major)
	if err == nil && majorVersion < 5 {

		if i.cfg.ArchWorker == "both" {
			return fmt.Errorf(
				"multiple worker installations " +
					"is only available from version 5")
		}

	}
	// ====================================================================== //

	steps := []ui.Step{

		// cryosparc_master
		i.backupCryosparcDatabaseStep(
			i.cfg.CryosparcPath,
			i.cfg.DbPath,
			filepath.Join("/data/user", os.Getenv("USER"), "cryosparc_backup")),

		i.CryosparcmStopStep(
			i.cfg.HostName, i.cfg.CryosparcPath, false),

		i.checkInstallDirStep(i.cfg.CryosparcPath, i.cfg.CryosparcPath),

		i.createDbPathStep(i.cfg.DbPath),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"master",
			i.cfg.ArchMaster),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath, i.cfg.Version, "master", i.cfg.ArchMaster),

		i.setCryosparcBasePortStep(
			i.cfg.HostName,
			i.cfg.BasePort,
			portStart,
			portEnd,
			portCount,
		),

		i.InstallMasterStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath, i.cfg.SSDPath,
			i.cfg.ArchMaster),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"),
			i.cfg.License,
			"master"),

		i.CryosparcmStartStep(
			i.cfg.HostName, i.cfg.CryosparcPath, false),

		i.CryosparcmCreateUserStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.User.Email,
			i.cfg.User.Username,
			i.cfg.User.FirstName,
			i.cfg.User.LastName,
			i.cfg.CryosparcmHelp),
	}

	// cryosparc_worker
	var archList []string

	if i.cfg.ArchWorker == "both" {
		archList = []string{"x86_64", "aarch64"}
	} else {
		archList = []string{i.cfg.ArchWorker}
	}

	for _, arch := range archList {

		steps = append(steps,

			i.downloadCryosparcStep(
				i.cfg.CryosparcPath, i.cfg.Version,
				i.cfg.License, "worker", arch),

			i.ExtractArchiveStep(
				i.cfg.CryosparcPath, i.cfg.Version, "worker", arch),

			i.InstallWorkerStep(
				i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, arch),

			i.replaceLicenseIDStep(
				filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
				i.cfg.License,
				"worker"),
		)

		// ================================================================== //
		if majorVersion >= 5 {
			steps[len(steps)-1] = i.RenameWorkerDirectoryStep(
				i.cfg.CryosparcPath,
				arch,
			)
		}
		// ================================================================== //

		// Create all default lanes
		steps = append(steps, i.LanesCreateDefaultSteps(
			i.cfg.CryosparcPath,
			i.cfg.Lanes.CachePath,
			i.cfg.Lanes.Memory,
			i.cfg.Lanes.Gpus,
			i.cfg.Lanes.CpusPerTask,
			arch)...)
	}

	// Install all default lanes
	steps = append(steps, i.LanesInstallSteps(
		i.cfg.HostName,
		i.cfg.CryosparcPath,
		i.cfg.Lanes.Info,
		i.cfg.Lanes.Script,
		i.cfg.Lanes.Name,
		true, // bool to install all lanes created (--all)
		i.cfg.CryosparcmHelp,
	)...)

	steps = append(steps, i.HttpUrlStep(i.cfg.HostName))

	return i.runSteps(steps...)
}

// install master
func (i *Command) RunInstallMaster() error {
	return i.runSteps(

		i.backupCryosparcDatabaseStep(
			i.cfg.CryosparcPath,
			i.cfg.DbPath,
			filepath.Join("/data/user", os.Getenv("USER"), "cryosparc_backup")),

		i.CryosparcmStopStep(
			i.cfg.HostName, i.cfg.CryosparcPath, false),

		i.checkInstallDirStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"),
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master")),

		i.createDbPathStep(i.cfg.DbPath),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"master",
			i.cfg.ArchMaster),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			"master",
			i.cfg.ArchMaster),

		i.setCryosparcBasePortStep(
			i.cfg.HostName,
			i.cfg.BasePort,
			portStart,
			portEnd,
			portCount,
		),

		i.InstallMasterStep(
			i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath,
			i.cfg.SSDPath, i.cfg.ArchMaster),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"),
			i.cfg.License,
			"master"),

		i.CryosparcmStartStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			false),

		i.CryosparcmCreateUserStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.User.Email,
			i.cfg.User.Username,
			i.cfg.User.FirstName,
			i.cfg.User.LastName,
			false),

		i.HttpUrlStep(i.cfg.HostName),
	)
}

// install worker
func (i *Command) RunInstallWorker() error {

	steps := []ui.Step{

		i.checkInstallDirStep(
			filepath.Join(
				i.cfg.CryosparcPath,
				fmt.Sprintf("cryosparc_worker_%s", i.cfg.ArchWorker),
			),
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
		),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"worker",
			i.cfg.ArchWorker),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.ArchWorker),

		i.InstallWorkerStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			i.cfg.ArchWorker,
		),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
			i.cfg.License,
			"worker"),

		i.RenameWorkerDirectoryStep(i.cfg.CryosparcPath, i.cfg.ArchWorker),
	}

	// ====================================================================== //
	// Only after version 5 that CryoSPARC allows to have different cryosparc
	// worker installations
	// ====================================================================== //
	version := strings.TrimSpace(i.cfg.Version)

	major, _, ok := strings.Cut(version, ".")
	if !ok {
		major = version
	}

	if majorVersion, err := strconv.Atoi(major); err == nil &&
		majorVersion < 5 {

		// do not create a worker directory with arch
		steps[0] =
			i.checkInstallDirStep(
				filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
				filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
			)

		// do not run rename worker directory
		steps = steps[:len(steps)-1]
	}
	// ====================================================================== //

	return i.runSteps(steps...)
}

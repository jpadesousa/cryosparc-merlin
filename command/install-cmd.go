package command

import (
	"cryosparc-merlin/ui"
	"os"
	"path/filepath"
)

// install complete
func (i *Command) RunInstallComplete() error {

	steps := []ui.Step{

		i.backupCryosparcDatabaseStep(
			i.cfg.CryosparcPath,
			i.cfg.DbPath,
			filepath.Join("/data/user", os.Getenv("USER"), "cryosparc_backup")),

		i.CryosparcmStopStep(
			i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),

		i.checkInstallDirStep(i.cfg.CryosparcPath),

		i.createDbPathStep(i.cfg.DbPath),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"master",
			i.cfg.ArchMaster),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath, i.cfg.Version, "master", i.cfg.ArchMaster),

		i.setCryosparcBasePortStep(i.cfg.HostName, i.cfg.BasePort),

		i.InstallMasterStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath, i.cfg.SSDPath,
			i.cfg.ArchMaster, i.cfg.BasePort),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"),
			i.cfg.License,
			"master"),

		i.CryosparcmStartStep(
			i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),

		i.CryosparcmCreateUserStep(
			i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath, i.cfg.Version,
			i.cfg.License, "worker", i.cfg.ArchWorker),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.ArchWorker),

		i.InstallWorkerStep(
			i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, i.cfg.ArchWorker),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
			i.cfg.License,
			"worker"),
	}

	// Create all default lanes
	steps = append(steps, i.LanesCreateDefaultSteps(
		i.cfg.CryosparcPath,
		i.cfg.Lanes.CachePath,
		i.cfg.Lanes.Memory,
		i.cfg.Lanes.Gpus,
		i.cfg.Lanes.CpusPerTask)...)

	// Install all default lanes
	steps = append(steps, i.LanesInstallSteps(
		i.cfg.HostName,
		i.cfg.CryosparcPath,
		i.cfg.Lanes.Info,
		i.cfg.Lanes.Script,
		i.cfg.Lanes.Name,
		true,  // bool to install all lanes created (--all)
		false, // bool to show help (--help)
	)...)

	// Print HTTP url
	steps = append(steps, i.HttpUrlStep(i.cfg.HostName, i.cfg.BasePort))

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
			i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),

		i.checkInstallDirStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master")),

		i.createDbPathStep(i.cfg.DbPath),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"master",
			i.cfg.Arch),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			"master",
			i.cfg.Arch),

		i.setCryosparcBasePortStep(i.cfg.HostName, i.cfg.BasePort),

		i.InstallMasterStep(
			i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath,
			i.cfg.SSDPath, i.cfg.Arch, i.cfg.BasePort),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"),
			i.cfg.License,
			"master"),

		i.CryosparcmStartStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.CryosparcmHelp),

		i.HttpUrlStep(i.cfg.HostName, i.cfg.BasePort),
	)
}

// install worker
func (i *Command) RunInstallWorker() error {
	return i.runSteps(

		i.checkInstallDirStep(
			filepath.Join(i.cfg.CryosparcPath,
				"cryosparc_worker")),

		i.downloadCryosparcStep(
			i.cfg.CryosparcPath,
			i.cfg.Version,
			i.cfg.License,
			"worker",
			i.cfg.Arch),

		i.ExtractArchiveStep(
			i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.Arch),

		i.InstallWorkerStep(
			i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, i.cfg.Arch),

		i.replaceLicenseIDStep(
			filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"),
			i.cfg.License,
			"worker"),
	)
}

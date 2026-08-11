package command

import (
	"cryosparc-merlin/ui"
	"os"
	"path/filepath"
)

// install complete
func (i *Command) RunInstallComplete() error {

	steps := []ui.Step{
		i.createCryosparcPathStep(i.cfg.CryosparcPath),
		i.backupCryosparcDatabaseStep(i.cfg.CryosparcPath, i.cfg.DbPath, filepath.Join("/data/user", os.Getenv("USER"), "cryosparc_backup")),
		i.createDbPathStep(i.cfg.DbPath),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "master", i.cfg.ArchMaster),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "master", i.cfg.ArchMaster),
		i.InstallMasterStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath, i.cfg.SSDPath, i.cfg.ArchMaster, i.cfg.BasePort),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"), i.cfg.License, "master"),
		i.CryosparcmStartStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),
		i.CryosparcmCreateUserStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "worker", i.cfg.ArchWorker),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.ArchWorker),
		i.InstallWorkerStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, i.cfg.ArchWorker),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"), i.cfg.License, "worker"),
	}

	// Create all default lanes
	steps = append(steps, i.RunLanesCreateDefaultSteps(
		i.cfg.CryosparcPath,
		i.cfg.Lanes.CachePath,
		i.cfg.Lanes.Memory,
		i.cfg.Lanes.Gpus,
		i.cfg.Lanes.CpusPerTask)...)

	// Install all default lanes
	steps = append(steps, i.RunLanesInstallSteps(
		i.cfg.HostName,
		i.cfg.CryosparcPath,
		i.cfg.Lanes.Info,
		i.cfg.Lanes.Script,
		i.cfg.Lanes.Name,
		true, // bool to install all lanes created
		i.cfg.CryosparcmHelp)...)

	return i.runSteps(steps...)
}

// install master
func (i *Command) RunInstallMaster() error {
	return i.runSteps(
		i.checkInstallDirStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_master")),
		i.backupCryosparcDatabaseStep(i.cfg.CryosparcPath, i.cfg.DbPath, filepath.Join("/data/user", os.Getenv("USER"), "cryosparc_backup")),
		i.createDbPathStep(i.cfg.DbPath),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "master", i.cfg.Arch),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "master", i.cfg.Arch),
		i.InstallMasterStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath, i.cfg.SSDPath, i.cfg.Arch, i.cfg.BasePort),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"), i.cfg.License, "master"),
		i.CryosparcmStartStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),
	)
}

// install worker
func (i *Command) RunInstallWorker() error {
	return i.runSteps(
		i.checkInstallDirStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker")),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "worker", i.cfg.Arch),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.Arch),
		i.InstallWorkerStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, i.cfg.Arch),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"), i.cfg.License, "worker"),
	)
}

package command

import (
	"path/filepath"
)

// install complete
func (i *Command) RunInstallComplete() error {
	return i.runSteps(
		i.createCryosparcPathStep(i.cfg.CryosparcPath),
		i.createDbPathStep(i.cfg.DbPath),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "master", i.cfg.ArchMaster),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "master", i.cfg.ArchMaster),
		i.InstallMasterStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License,
			i.cfg.HostName, i.cfg.DbPath, i.cfg.SSDPath, i.cfg.ArchMaster, i.cfg.BasePort),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_master"), i.cfg.License, "master"),
		i.CryosparcmStartStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),
		i.downloadCryosparcStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, "worker", i.cfg.ArchWorker),
		i.ExtractArchiveStep(i.cfg.CryosparcPath, i.cfg.Version, "worker", i.cfg.ArchWorker),
		i.InstallWorkerStep(i.cfg.CryosparcPath, i.cfg.Version, i.cfg.License, i.cfg.ArchWorker),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_worker"), i.cfg.License, "worker"),
	)
}

// install master
func (i *Command) RunInstallMaster() error {
	return i.runSteps(
		i.checkInstallDirStep(filepath.Join(i.cfg.CryosparcPath, "cryosparc_master")),
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

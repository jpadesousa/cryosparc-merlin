package command

import (
	"path/filepath"
)

// install complete
func (i *Command) RunInstallComplete() error {
	return i.runSteps(
		i.createHomeDirStep(i.cfg.HomeDir),
		i.createDbDirStep(i.cfg.DbDir),
		i.downloadCryosparcStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, "master", i.cfg.ArchMaster),
		i.ExtractArchiveStep(i.cfg.HomeDir, i.cfg.Version, "master", i.cfg.ArchMaster),
		i.InstallMasterStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License,
			i.cfg.RemoteHost, i.cfg.DbDir, i.cfg.ArchMaster, i.cfg.BasePort),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.HomeDir, "cryosparc_master"), i.cfg.License, "master"),
		i.CryosparcmStartStep(i.cfg.RemoteHost, i.cfg.HomeDir),
		i.downloadCryosparcStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, "worker", i.cfg.ArchWorker),
		i.ExtractArchiveStep(i.cfg.HomeDir, i.cfg.Version, "worker", i.cfg.ArchWorker),
		i.InstallWorkerStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, i.cfg.ArchWorker),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.HomeDir, "cryosparc_worker"), i.cfg.License, "worker"),
	)
}

// install master
func (i *Command) RunInstallMaster() error {
	return i.runSteps(
		i.checkInstallDirStep(filepath.Join(i.cfg.HomeDir, "cryosparc_master")),
		i.createDbDirStep(i.cfg.DbDir),
		i.downloadCryosparcStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, "master", i.cfg.Arch),
		i.ExtractArchiveStep(i.cfg.HomeDir, i.cfg.Version, "master", i.cfg.Arch),
		i.InstallMasterStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License,
			i.cfg.RemoteHost, i.cfg.DbDir, i.cfg.Arch, i.cfg.BasePort),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.HomeDir, "cryosparc_master"), i.cfg.License, "master"),
		i.CryosparcmStartStep(i.cfg.RemoteHost, i.cfg.HomeDir),
	)
}

// install worker
func (i *Command) RunInstallWorker() error {
	return i.runSteps(
		i.checkInstallDirStep(filepath.Join(i.cfg.HomeDir, "cryosparc_worker")),
		i.downloadCryosparcStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, "worker", i.cfg.Arch),
		i.ExtractArchiveStep(i.cfg.HomeDir, i.cfg.Version, "worker", i.cfg.Arch),
		i.InstallWorkerStep(i.cfg.HomeDir, i.cfg.Version, i.cfg.License, i.cfg.Arch),
		i.replaceLicenseIDStep(filepath.Join(i.cfg.HomeDir, "cryosparc_worker"), i.cfg.License, "worker"),
	)
}

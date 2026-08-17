package command

func (i *Command) RunLanesCreateCmd() error {
	return i.runSteps(
		i.LanesCreateStep(
			i.cfg.CryosparcPath,
			i.cfg.Lanes.Name,
			i.cfg.Lanes.CachePath,
			i.cfg.Lanes.Memory,
			i.cfg.Lanes.Time,
			i.cfg.Lanes.Partition,
			i.cfg.Lanes.Gpus,
			i.cfg.Lanes.CpusPerTask,
			i.cfg.Lanes.Cluster,
			i.cfg.ArchWorker,
		),
	)
}

func (i *Command) RunLanesCreateDefaultCmd() error {
	return i.runSteps(i.LanesCreateDefaultSteps(
		i.cfg.CryosparcPath,
		i.cfg.Lanes.CachePath,
		i.cfg.Lanes.Memory,
		i.cfg.Lanes.Gpus,
		i.cfg.Lanes.CpusPerTask,
		i.cfg.ArchWorker)...)
}

func (i *Command) RunLanesInstallCmd() error {
	return i.runSteps(i.LanesInstallSteps(
		i.cfg.HostName,
		i.cfg.CryosparcPath,
		i.cfg.Lanes.Info,
		i.cfg.Lanes.Script,
		i.cfg.Lanes.Name,
		i.cfg.Lanes.InstallAll,
		i.cfg.CryosparcmHelp)...)
}

func (i *Command) RunLanesRemoveCmd() error {
	return i.runSteps(
		i.LanesRemoveStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.Lanes.Name,
			i.cfg.CryosparcmHelp,
		),
	)
}

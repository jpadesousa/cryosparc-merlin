package command

func (i *Command) RunLanesCreateCmd() error {
	return i.runSteps(
		i.RunLanesCreateStep(
			i.cfg.CryosparcPath,
			i.cfg.LanesCreate.Name,
			i.cfg.LanesCreate.CachePath,
			i.cfg.LanesCreate.Memory,
			i.cfg.LanesCreate.Time,
			i.cfg.LanesCreate.Partition,
			i.cfg.LanesCreate.Gpus,
			i.cfg.LanesCreate.CpusPerTask,
			i.cfg.LanesCreate.Cluster,
		),
	)
}

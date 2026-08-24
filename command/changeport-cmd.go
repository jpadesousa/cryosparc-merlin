package command

func (i *Command) RunChangePort() error {
	return i.runSteps(
		i.CryosparcmChangePortStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.BasePort,
			portStart,
			portEnd,
			portCount,
			i.cfg.CryosparcmHelp,
		),
	)
}

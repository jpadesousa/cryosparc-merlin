package command

func (i *Command) RunCryosparcm(args []string) error {
	return i.runSteps(
		i.CryosparcmStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp, args...),
	)
}

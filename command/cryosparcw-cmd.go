package command

func (i *Command) RunCryosparcw(args []string) error {
	return i.runSteps(
		i.CryosparcwStep(
			i.cfg.CryosparcPath,
			i.cfg.ArchWorker,
			i.cfg.CryosparcwHelp,
			args...),
	)
}

package command

func (i *Command) RunUserCreateCmd() error {
	return i.runSteps(
		i.CryosparcmCreateUserStep(
			i.cfg.HostName,
			i.cfg.CryosparcPath,
			i.cfg.User.Email,
			i.cfg.User.Username,
			i.cfg.User.FirstName,
			i.cfg.User.LastName,
			i.cfg.CryosparcmHelp,
		),
	)
}

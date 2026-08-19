package command

func (i *Command) RunCheckInstances() error {
	return i.runSteps(
		i.CheckInstancesStep(
			i.cfg.HostName,
		),
	)
}

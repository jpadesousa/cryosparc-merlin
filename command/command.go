package command

import (
	"cryosparc-merlin/config"
	"cryosparc-merlin/ui"
)

type Command struct {
	cfg config.Config
}

func New(cfg config.Config) *Command {
	return &Command{cfg: cfg}
}

func (i *Command) runSteps(steps ...ui.Step) error {
	return ui.Run(steps)
}

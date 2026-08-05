package command

// cryosparcm status
func (i *Command) RunCryosparcmStatus() error {
	return i.runSteps(
		i.CryosparcmStatusStep(i.cfg.RemoteHost, i.cfg.HomeDir),
	)
}

// cryosparcm start
func (i *Command) RunCryosparcmStart() error {
	return i.runSteps(
		i.CryosparcmStartStep(i.cfg.RemoteHost, i.cfg.HomeDir),
	)
}

// cryosparcm stop
func (i *Command) RunCryosparcmStop() error {
	return i.runSteps(
		i.CryosparcmStopStep(i.cfg.RemoteHost, i.cfg.HomeDir),
	)
}

// cryosparcm restart
func (i *Command) RunCryosparcmRestart() error {
	return i.runSteps(
		i.CryosparcmRestartStep(i.cfg.RemoteHost, i.cfg.HomeDir),
	)
}

// cryosparcm createuser
func (i *Command) RunCryosparcmCreateUser() error {
	return i.runSteps(
		i.CryosparcmCreateUserStep(
			i.cfg.RemoteHost,
			i.cfg.HomeDir,
			i.cfg.Email,
			i.cfg.Username,
			i.cfg.FirstName,
			i.cfg.LastName,
		),
	)
}

// cryosparcm resetpassword
func (i *Command) RunCryosparcmResetPassword() error {
	return i.runSteps(
		i.CryosparcmResetPasswordStep(
			i.cfg.RemoteHost,
			i.cfg.HomeDir,
			i.cfg.Email,
		),
	)
}

// cryosparcm update
func (i *Command) RunCryosparcmUpdate() error {
	return i.runSteps(
		i.CryosparcmUpdateStep(
			i.cfg.RemoteHost,
			i.cfg.HomeDir,
			i.cfg.CryosparcmUpdate.Version,
			i.cfg.CryosparcmUpdate.Check,
			i.cfg.CryosparcmUpdate.List,
			i.cfg.CryosparcmUpdate.Override,
			i.cfg.CryosparcmUpdate.DownloadOnly,
			i.cfg.CryosparcmUpdate.SkipDownload,
		),
	)
}

// cryosparcm patch
func (i *Command) RunCryosparcmPatch() error {
	return i.runSteps(
		i.CryosparcmPatchStep(
			i.cfg.RemoteHost,
			i.cfg.HomeDir,
			i.cfg.CryosparcmPatch.Install,
			i.cfg.CryosparcmPatch.Download,
			i.cfg.CryosparcmPatch.Check,
			i.cfg.CryosparcmPatch.Yes,
			i.cfg.CryosparcmPatch.Force,
		),
	)
}

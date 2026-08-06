package command

func (i *Command) RunCryosparcm(args []string) error {
	return i.runSteps(
		i.CryosparcmStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp, args...),
	)
}

// cryosparcm status
// func (i *Command) RunCryosparcmStatus() error {
// 	return i.runSteps(
// 		i.CryosparcmStatusStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.CryosparcmHelp),
// 	)
// }

// cryosparcm start
// func (i *Command) RunCryosparcmStart() error {
// 	return i.runSteps(
// 		i.CryosparcmStartStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.Help),
// 	)
// }

// // cryosparcm stop
// func (i *Command) RunCryosparcmStop() error {
// 	return i.runSteps(
// 		i.CryosparcmStopStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.Help),
// 	)
// }

// // cryosparcm restart
// func (i *Command) RunCryosparcmRestart() error {
// 	return i.runSteps(
// 		i.CryosparcmRestartStep(i.cfg.HostName, i.cfg.CryosparcPath, i.cfg.Help),
// 	)
// }

// // cryosparcm createuser
// func (i *Command) RunCryosparcmCreateUser() error {
// 	return i.runSteps(
// 		i.CryosparcmCreateUserStep(
// 			i.cfg.HostName,
// 			i.cfg.CryosparcPath,
// 			i.cfg.CryosparcmUser.Email,
// 			i.cfg.CryosparcmUser.Username,
// 			i.cfg.CryosparcmUser.FirstName,
// 			i.cfg.CryosparcmUser.LastName,
// 			i.cfg.Help,
// 		),
// 	)
// }

// // cryosparcm resetpassword
// func (i *Command) RunCryosparcmResetPassword() error {
// 	return i.runSteps(
// 		i.CryosparcmResetPasswordStep(
// 			i.cfg.HostName,
// 			i.cfg.CryosparcPath,
// 			i.cfg.CryosparcmUser.Email,
// 			i.cfg.Help,
// 		),
// 	)
// }

// // cryosparcm update
// func (i *Command) RunCryosparcmUpdate() error {
// 	return i.runSteps(
// 		i.CryosparcmUpdateStep(
// 			i.cfg.HostName,
// 			i.cfg.CryosparcPath,
// 			i.cfg.CryosparcmUpdate.Version,
// 			i.cfg.CryosparcmUpdate.Check,
// 			i.cfg.CryosparcmUpdate.List,
// 			i.cfg.CryosparcmUpdate.Override,
// 			i.cfg.CryosparcmUpdate.DownloadOnly,
// 			i.cfg.CryosparcmUpdate.SkipDownload,
// 			i.cfg.Help,
// 		),
// 	)
// }

// // cryosparcm patch
// func (i *Command) RunCryosparcmPatch() error {
// 	return i.runSteps(
// 		i.CryosparcmPatchStep(
// 			i.cfg.HostName,
// 			i.cfg.CryosparcPath,
// 			i.cfg.CryosparcmPatch.Install,
// 			i.cfg.CryosparcmPatch.Download,
// 			i.cfg.CryosparcmPatch.Check,
// 			i.cfg.CryosparcmPatch.Yes,
// 			i.cfg.CryosparcmPatch.Force,
// 			i.cfg.Help,
// 		),
// 	)
// }

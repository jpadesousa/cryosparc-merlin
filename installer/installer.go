package installer

import (
	"cryosparc-install/config"
	"errors"
	"fmt"
	"os"
	"time"

	"cryosparc-install/ui"

	"charm.land/huh/v2"
)

type Installer struct {
	cfg config.Config
}

func New(cfg config.Config) *Installer {
	return &Installer{cfg: cfg}
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// func directoryExistsAndNotEmpty(path string) (bool, error) {
// 	info, err := os.Stat(path)
// 	switch {
// 	case errors.Is(err, os.ErrNotExist):
// 		return false, nil

// 	case err != nil:
// 		return false, err

// 	case !info.IsDir():
// 		return false, fmt.Errorf("%q exists but is not a directory", path)
// 	}

// 	f, err := os.Open(path)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer f.Close()

// 	_, err = f.Readdirnames(1)
// 	if errors.Is(err, io.EOF) {
// 		return false, nil // directory is empty
// 	}
// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil // directory exists and contains at least one entry
// }

var overwrite bool
var action string
var ErrInstallationCancelled = errors.New("installation cancelled by user")

func (i *Installer) Run() error {
	steps := []ui.Step{
		{
			Message: fmt.Sprintf("Creating home directory %s", i.cfg.HomeDir),
			Condition: func() bool {
				time.Sleep(500 * time.Millisecond)
				return directoryExists(i.cfg.HomeDir)
			},
			Prompt: func() *huh.Form {
				return huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title(fmt.Sprintf(
								"%q already exists",
								i.cfg.HomeDir,
							)).
							Options(
								huh.NewOption("Overwrite installation", "overwrite"),
								huh.NewOption("Cancel installation", "cancel"),
							).
							Value(&action),
					),
				)
			},
			// SkipMessage: "Home directory already exists",
			Action: func(update func(float64)) error {
				if directoryExists(i.cfg.HomeDir) {
					if action == "cancel" {
						return ErrInstallationCancelled
					}

					if err := os.RemoveAll(i.cfg.HomeDir); err != nil {
						return err
					}
				}

				return i.createHomeDirectory(update)
			},
		},
		// {
		// 	Message: fmt.Sprintf("Creating database directory %s", i.cfg.DbDir),
		// 	Action:  i.createDatabaseDirectory,
		// },
		{
			Message: fmt.Sprintf("Creating database directory %s", i.cfg.DbDir),
			Condition: func() bool {
				time.Sleep(500 * time.Millisecond)
				return directoryExists(i.cfg.DbDir)
			},
			SkipMessage: "Database directory already exists",
			Action:      i.createDatabaseDirectory,
		},
		{
			Message:         fmt.Sprintf("Downloading CryoSPARC master v%s", i.cfg.Release),
			ShowProgressBar: true,
			Action:          i.downloadMaster,
		},
		{
			Message: fmt.Sprintf("Extracting CryoSPARC master v%s", i.cfg.Release),
			Action: func(update func(float64)) error {
				return i.extractArchive(update, "cryosparc_master.tar.gz")
			},
		},
		{
			Message: fmt.Sprintf("Installing CryoSPARC master v%s", i.cfg.Release),
			Exec:    i.installMaster,
		},
	}

	return ui.Run(steps)
}

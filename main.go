package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cryosparc-merlin-install/config"
	"cryosparc-merlin-install/installer"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/spf13/cobra"
)

// Color and style variables
var (
	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14"))

	helpTitleStyle = lipgloss.NewStyle().
			Foreground(charmtone.Charple).Bold(true)
)

// Default flags variables
var (
	user              = os.Getenv("USER")
	defaultHomeDir    = filepath.Join("/data/user", user, "cryosparc")
	defaultDbDir      = filepath.Join(defaultHomeDir, "database")
	defaultRelease    = "4.7.1"
	defaultRemoteHost = "service03.merlin7.psi.ch"
)

var cfg config.Config

func main() {

	// CLI Create a new cobra Command
	cobraCmd := &cobra.Command{
		Use: "cryosparc-merlin-install",
		Long: fmt.Sprintf(`%s

Install and configure a CryoSPARC instance on the Merlin 7 cluster.`,
			helpTitleStyle.Render("DESCRIPTION")),
		RunE: runApp,
	}

	// =============================================
	// CLI Flags
	// =============================================

	flags := cobraCmd.Flags()

	flags.StringVar(&cfg.Release, "release", defaultRelease,
		"CryoSPARC release to install",
	)

	flags.StringVar(&cfg.HomeDir, "homedir", defaultHomeDir,
		"Installation directory for CryoSPARC",
	)

	flags.StringVar(&cfg.DbDir, "dbdir", defaultDbDir,
		"Directory for the CryoSPARC database",
	)

	flags.StringVar(&cfg.RemoteHost, "remotehost", defaultRemoteHost,
		"Hostname used to access the CryoSPARC web interface",
	)

	flags.StringVar(&cfg.License, "license", "",
		"CryoSPARC license ID",
	)

	flags.UintVar(&cfg.BasePort, "port", 0,
		"Base TCP port for CryoSPARC services",
	)

	// Required flags
	_ = cobraCmd.MarkFlagRequired("license")

	// Check flags
	// cobraCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
	// 	if cfg.basePort > 0 {
	// 		return fmt.Errorf("--base-port must be greater than zero")
	// 	}

	// 	return nil
	// }

	// =============================================
	// Execute app
	// =============================================

	// Set new flag display settings
	cobraCmd.Flags().SetInterspersed(false)
	cobraCmd.Flags().SortFlags = false

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Execute cobra command
	if err := fang.Execute(context.Background(), cobraCmd); err != nil {
		fmt.Fprintf(os.Stderr, "  Run '%s --help' for usage.\n\n", cobraCmd.CommandPath())
		os.Exit(1)
	}
}

func runApp(cmd *cobra.Command, _ []string) error {

	// If it is a beta release and the user did not change homedir, homedir
	// will be set to the user's home directory + /cryosparc_beta
	if strings.Contains(cfg.Release, "beta") && !cmd.Flags().Changed("homedir") {
		cfg.HomeDir = filepath.Join("/data/user", user, "cryosparc_beta")
	}

	// If the user did not change database directory and homedir
	// is the different than default, the dbdir will be set inside the homedir
	if !cmd.Flags().Changed("dbdir") && cfg.HomeDir != defaultHomeDir {
		cfg.DbDir = filepath.Join(cfg.HomeDir, "database")
	}

	inst := installer.New(cfg)
	inst.Run()

	return nil
}

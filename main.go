package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cryosparc-merlin/command"
	"cryosparc-merlin/config"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/spf13/cobra"
)

// Color and style variables
var (
	helpTitleStyle = lipgloss.NewStyle().
		Foreground(charmtone.Charple).Bold(true)
)

// Default flags variables
var (
	user                 = os.Getenv("USER")
	defaultCryosparcPath = filepath.Join("/data/user", user, "cryosparc")
	defaultDbPath        = filepath.Join(defaultCryosparcPath, "database")
	defaultSSDPath       = "/scratch"
	defaultVersion       = "4.7.1"
	defaultHostName      = "service03.merlin7.psi.ch"
	defaultArchMaster    = "x86_64"
	defaultArchWorker    = "x86_64"
	defaultArch          = "x86_64"
)

var cfg config.Config

// Valid architectures
var validArch = map[string]struct{}{
	"x86_64":  {},
	"aarch64": {},
}

func main() {

	// =========================================================================
	// CLI Flags
	// =========================================================================

	// Base command
	cobraCmd := newCommand(
		"cryosparc-merlin",
		"CryoSPARC helper for the Merlin cluster",
		runHelpCmd,
	)

	// cryosparcm commands
	cryosparcmCmd := newCommand(
		"cryosparcm",
		"Cryosparcm commands",
		runCmdArgs((*command.Command).RunCryosparcm),
	)

	// Install commands
	installCmd := newCommand(
		"install",
		"Installation commands",
		runHelpCmd,
	)

	installCompleteCmd := newCommand(
		"complete",
		"Complete CryoSPARC installation",
		runCmd((*command.Command).RunInstallComplete),
	)

	installCompleteCmd.PreRunE = func(cmd *cobra.Command, args []string) error {

		if _, ok := validArch[cfg.ArchMaster]; !ok {
			return fmt.Errorf(
				"invalid value for --arch-master: %q (must be one of: x86_64, aarch64)",
				cfg.ArchMaster,
			)
		}

		if _, ok := validArch[cfg.ArchWorker]; !ok {
			return fmt.Errorf(
				"invalid value for --arch-worker: %q (must be one of: x86_64, aarch64)",
				cfg.ArchWorker,
			)
		}

		return nil
	}

	installMasterCmd := newCommand(
		"master",
		"Install CryoSPARC master",
		runCmd((*command.Command).RunInstallMaster),
	)

	installMasterCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if _, ok := validArch[cfg.Arch]; !ok {
			return fmt.Errorf(
				"invalid value for --arch: %q (must be one of: x86_64, aarch64)",
				cfg.Arch,
			)
		}
		return nil
	}

	installWorkerCmd := newCommand(
		"worker",
		"Install CryoSPARC worker",
		runCmd((*command.Command).RunInstallWorker),
	)

	installWorkerCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if _, ok := validArch[cfg.Arch]; !ok {
			return fmt.Errorf(
				"invalid value for --arch: %q (must be one of: x86_64, aarch64)",
				cfg.Arch,
			)
		}
		return nil
	}

	installCmd.AddCommand(
		installCompleteCmd,
		installMasterCmd,
		installWorkerCmd,
	)

	// Base commands
	cobraCmd.AddCommand(cryosparcmCmd, installCmd)

	// =============================================
	// cryosparcm CLI Flags
	// =============================================

	addCryosparcmFlags(cryosparcmCmd)

	// =============================================
	// Install CLI Flags
	// =============================================

	// complete
	addCommonInstallFlags(installCompleteCmd)
	addInstallCompleteFlags(installCompleteCmd)

	// master
	addCommonInstallFlags(installMasterCmd)
	addInstallMasterFlags(installMasterCmd)

	// worker
	addCommonInstallFlags(installWorkerCmd)
	addInstallWorkerFlags(installWorkerCmd)

	// =============================================
	// Execute app
	// =============================================

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Execute cobra command
	if err := fang.Execute(context.Background(), cobraCmd); err != nil {
		fmt.Fprintf(os.Stderr, "  Run '%s --help' for usage.\n\n", cobraCmd.CommandPath())
		os.Exit(1)
	}
}

// =============================================================================
// Create new commands
// =============================================================================
func newCommand(
	use, desc string,
	runE func(*cobra.Command, []string) error,
) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: desc,
		Args:  cobra.ArbitraryArgs,
		Long: fmt.Sprintf(`%s

%s`,
			helpTitleStyle.Render("DESCRIPTION"),
			desc,
		),
		RunE: runE,
	}
}

// =============================================================================
// Prepare config
// =============================================================================
func prepareConfig(cmd *cobra.Command) {

	if !cmd.Flags().Changed("dbpath") &&
		cfg.CryosparcPath != defaultCryosparcPath {
		cfg.DbPath = filepath.Join(cfg.CryosparcPath, "database")
	}

}

// =============================================================================
// Run Commands
// =============================================================================

func runHelpCmd(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func runCmd(fn func(*command.Command) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		prepareConfig(cmd)
		inst := command.New(cfg)
		return fn(inst)
	}
}

// If there is a command that allows the user to fill with arguments without
// being predefined
func runCmdArgs(fn func(*command.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		prepareConfig(cmd)
		inst := command.New(cfg)
		return fn(inst, args)
	}
}

// =============================================================================
// Flags helper functions
// =============================================================================

// cryosparcm
func addCryosparcmFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.CryosparcPath, "cryosparc-path", defaultCryosparcPath,
		"Installation directory for CryoSPARC",
	)

	flags.StringVar(&cfg.HostName, "hostname", defaultHostName,
		"Hostname used to access the CryoSPARC web interface",
	)

	flags.BoolVar(&cfg.CryosparcmHelp, "list-commands", false,
		"List cryosparcm commands (only available from version 5)",
	)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

// Install
func addCommonInstallFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.License, "license", "",
		"CryoSPARC license ID",
	)
	flags.StringVar(&cfg.Version, "version", defaultVersion,
		"CryoSPARC version to install",
	)
	flags.StringVar(&cfg.CryosparcPath, "cryosparc-path", defaultCryosparcPath,
		"Installation directory for CryoSPARC",
	)

	_ = cmd.MarkFlagRequired("license")

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addInstallMasterFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.DbPath, "dbpath", defaultDbPath,
		"CryoSPARC database directory",
	)
	flags.StringVar(&cfg.SSDPath, "ssdpath", defaultSSDPath,
		"SSD cache directory",
	)
	flags.StringVar(&cfg.HostName, "hostname", defaultHostName,
		"Hostname used to access the CryoSPARC web interface",
	)
	flags.UintVar(&cfg.BasePort, "port", 0,
		"Base TCP port for CryoSPARC services",
	)

	flags.StringVar(&cfg.Arch, "arch", defaultArch,
		"Master architecture (x86_64 or aarch64)",
	)
}

func addInstallWorkerFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.Arch, "arch", defaultArch,
		"Worker architecture (x86_64 or aarch64)",
	)
}

func addInstallCompleteFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.DbPath, "dbpath", defaultDbPath,
		"CryoSPARC database directory",
	)
	flags.StringVar(&cfg.HostName, "hostname", defaultHostName,
		"Hostname used to access the CryoSPARC web interface",
	)
	flags.UintVar(&cfg.BasePort, "port", 0,
		"Base TCP port for CryoSPARC services",
	)

	flags.StringVar(&cfg.ArchMaster, "arch-master", defaultArchMaster,
		"Master architecture (x86_64 or aarch64)",
	)

	flags.StringVar(&cfg.ArchWorker, "arch-worker", defaultArchWorker,
		"Worker architecture (x86_64 or aarch64)",
	)
}

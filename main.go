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
	pathStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	helpTitleStyle = lipgloss.NewStyle().
			Foreground(charmtone.Charple).Bold(true)
)

// Default flags variables
var (
	user                 = os.Getenv("USER")
	defaultHomeDir       = filepath.Join("/data/user", user, "cryosparc")
	defaultCryosparcmDir = filepath.Join(defaultHomeDir, "cryosparc_master")
	defaultDbDir         = filepath.Join(defaultHomeDir, "database")
	defaultVersion       = "4.7.1"
	defaultRemoteHost    = "service03.merlin7.psi.ch"
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
		runHelpCmd,
	)

	cryosparcmStatusCmd := newCommand(
		"status",
		"Check CryoSPARC status",
		runCmd((*command.Command).RunCryosparcmStatus),
	)

	cryosparcmStartCmd := newCommand(
		"start",
		"Start CryoSPARC",
		runCmd((*command.Command).RunCryosparcmStart),
	)

	cryosparcmStopCmd := newCommand(
		"stop",
		"Stop CryoSPARC",
		runCmd((*command.Command).RunCryosparcmStop),
	)

	cryosparcmRestartCmd := newCommand(
		"restart",
		"Restart CryoSPARC",
		runCmd((*command.Command).RunCryosparcmRestart),
	)

	cryosparcmCreateUserCmd := newCommand(
		"createuser",
		"Create CryoSPARC user",
		runCmd((*command.Command).RunCryosparcmCreateUser),
	)

	cryosparcmResetPasswordCmd := newCommand(
		"resetpassword",
		"Reset CryoSPARC password",
		runCmd((*command.Command).RunCryosparcmResetPassword),
	)

	cryosparcmUpdateCmd := newCommand(
		"update",
		"Update CryoSPARC",
		runCmd((*command.Command).RunCryosparcmUpdate),
	)

	cryosparcmPatchCmd := newCommand(
		"patch",
		"Patch CryoSPARC",
		runCmd((*command.Command).RunCryosparcmPatch),
	)

	// cryosparcmUpdateCmd := newCommand(
	// 	"update",
	// 	"Update CryoSPARC instance",
	// 	runCmdArgs((*command.Command).RunCryosparcmUpdate),
	// )

	// cryosparcmUpdateCmd.DisableFlagParsing = true

	cryosparcmCmd.AddCommand(
		cryosparcmStatusCmd,
		cryosparcmStartCmd,
		cryosparcmStopCmd,
		cryosparcmRestartCmd,
		cryosparcmCreateUserCmd,
		cryosparcmResetPasswordCmd,
		cryosparcmUpdateCmd,
		cryosparcmPatchCmd,
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

	installCompleteArmCmd := newCommand(
		"complete_arm",
		"Complete CryoSPARC installation to run on GH nodes",
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
		installCompleteArmCmd,
		installMasterCmd,
		installWorkerCmd,
	)

	// Base commands
	cobraCmd.AddCommand(cryosparcmCmd, installCmd)

	// =============================================
	// cryosparcm CLI Flags
	// =============================================

	// status
	addCryosparcmCommonFlags(cryosparcmStatusCmd)

	// start
	addCryosparcmCommonFlags(cryosparcmStartCmd)

	// stop
	addCryosparcmCommonFlags(cryosparcmStopCmd)

	// restart
	addCryosparcmCommonFlags(cryosparcmRestartCmd)

	// createuser
	addCryosparcmCommonFlags(cryosparcmCreateUserCmd)
	addCryosparcmCreateUserFlags(cryosparcmCreateUserCmd)

	// resetpassword
	addCryosparcmCommonFlags(cryosparcmResetPasswordCmd)
	addCryosparcmPasswordFlags(cryosparcmResetPasswordCmd)

	// update
	addCryosparcmCommonFlags(cryosparcmUpdateCmd)
	addCryosparcmUpdateFlags(cryosparcmUpdateCmd)

	// patch
	addCryosparcmCommonFlags(cryosparcmPatchCmd)
	addCryosparcmPatchFlags(cryosparcmPatchCmd)

	// =============================================
	// Install CLI Flags
	// =============================================

	// complete
	addCommonInstallFlags(installCompleteCmd)
	addInstallCompleteFlags(installCompleteCmd)

	// complete_arm
	addInstallCompletArmFlags(installCompleteArmCmd)

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
	// if strings.Contains(cfg.Version, "beta") &&
	// 	!cmd.Flags().Changed("homedir") {
	// 	cfg.HomeDir = filepath.Join("/data/user", user, "cryosparc_beta")
	// }

	if !cmd.Flags().Changed("dbdir") &&
		cfg.HomeDir != defaultHomeDir {
		cfg.DbDir = filepath.Join(cfg.HomeDir, "database")
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
func addCryosparcmCommonFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.HomeDir, "homedir", defaultHomeDir,
		"Installation directory for CryoSPARC",
	)

	flags.StringVar(&cfg.RemoteHost, "remotehost", defaultRemoteHost,
		"Hostname used to access the CryoSPARC web interface",
	)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addCryosparcmPasswordFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.Email, "email", "",
		"CryoSPARC user email",
	)
}

func addCryosparcmCreateUserFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.Username, "username", "",
		"CryoSPARC user username",
	)

	flags.StringVar(&cfg.FirstName, "firstname", "",
		"CryoSPARC user first name",
	)

	flags.StringVar(&cfg.LastName, "lastname", "",
		"CryoSPARC user last name",
	)
}

func addCryosparcmUpdateFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.CryosparcmUpdate.Version, "version", "",
		"Specify the version of CryoSPARC to update to (e.g., --version=v4.0.0)",
	)

	flags.BoolVar(&cfg.CryosparcmUpdate.Check, "check", false,
		"Check if there are any updates available for CryoSPARC",
	)

	flags.BoolVar(&cfg.CryosparcmUpdate.List, "list", false,
		"List all available versions CryoSPARC can update to",
	)

	flags.BoolVar(&cfg.CryosparcmUpdate.Override, "override", false,
		"Update to the latest version of CryoSPARC, regardless of the version the instance is currently on",
	)

	flags.BoolVar(&cfg.CryosparcmUpdate.DownloadOnly, "download-only", false,
		"Download the master and worker update packages without updating",
	)

	flags.BoolVar(&cfg.CryosparcmUpdate.SkipDownload, "skip-download", false,
		"Update CryoSPARC with previously-downloaded master and worker packages",
	)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addCryosparcmPatchFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.BoolVar(&cfg.CryosparcmPatch.Install, "install", false,
		"manually install a downloaded patch file",
	)

	flags.BoolVar(&cfg.CryosparcmPatch.Download, "download", false,
		"download master and worker patches for manual installation",
	)

	flags.BoolVar(&cfg.CryosparcmPatch.Check, "check", false,
		"check to see if a patch is available",
	)

	flags.BoolVarP(&cfg.CryosparcmPatch.Yes, "yes", "y", false,
		"confirm patch installation without prompt",
	)

	flags.BoolVarP(&cfg.CryosparcmPatch.Force, "force", "f", false,
		"install latest patch again even if already installed",
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
	flags.StringVar(&cfg.HomeDir, "homedir", defaultHomeDir,
		"Installation directory for CryoSPARC",
	)

	_ = cmd.MarkFlagRequired("license")

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addInstallMasterFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.DbDir, "dbdir", defaultDbDir,
		"Directory for the CryoSPARC database",
	)
	flags.StringVar(&cfg.RemoteHost, "remotehost", defaultRemoteHost,
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

	flags.StringVar(&cfg.DbDir, "dbdir", defaultDbDir,
		"Directory for the CryoSPARC database",
	)
	flags.StringVar(&cfg.RemoteHost, "remotehost", defaultRemoteHost,
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

func addInstallCompletArmFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.License, "license", "",
		"CryoSPARC license ID",
	)
	flags.StringVar(&cfg.Version, "version", "5.1.0-privatebeta.2",
		"CryoSPARC version to install",
	)
	flags.StringVar(&cfg.HomeDir, "homedir", filepath.Join("/data/user", user, "cryosparc_beta"),
		"Installation directory for CryoSPARC",
	)

	flags.StringVar(&cfg.DbDir, "dbdir", filepath.Join("/data/user", user, "cryosparc_beta", "database"),
		"Directory for the CryoSPARC database",
	)
	flags.StringVar(&cfg.RemoteHost, "remotehost", defaultRemoteHost,
		"Hostname used to access the CryoSPARC web interface",
	)
	flags.UintVar(&cfg.BasePort, "port", 0,
		"Base TCP port for CryoSPARC services",
	)

	flags.StringVar(&cfg.ArchMaster, "arch-master", "x86_64",
		"Master architecture (x86_64 or aarch64)",
	)

	flags.StringVar(&cfg.ArchWorker, "arch-worker", "aarch64",
		"Worker architecture (x86_64 or aarch64)",
	)

	_ = cmd.MarkFlagRequired("license")

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

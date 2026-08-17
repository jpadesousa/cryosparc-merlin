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

// Set version
var version = "dev"

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

	// Lanes
	defaultLaneMemory      = "ram_gb * 2"
	defaultLaneGpus        = "num_gpu"
	defaultLaneCpusPerTask = "num_cpu"
	defaultLaneCachePath   = "/scratch"
)

// Flag variables
var cfg config.Config

// Valid architectures
var validArch = map[string]struct{}{
	"x86_64":  {},
	"aarch64": {},
}

var validArchWorkerInstall = map[string]struct{}{
	"x86_64":  {},
	"aarch64": {},
	"both":    {},
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

	// cryosparcw commands
	cryosparcwCmd := newCommand(
		"cryosparcw",
		"Cryosparcw commands",
		runCmdArgs((*command.Command).RunCryosparcw),
	)

	cryosparcwCmd.PreRunE = preRunCmd("worker", false)

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

	installCompleteCmd.PreRunE = preRunCmd("both", true)

	installMasterCmd := newCommand(
		"master",
		"Install CryoSPARC master",
		runCmd((*command.Command).RunInstallMaster),
	)

	installMasterCmd.PreRunE = preRunCmd("master", true)

	installWorkerCmd := newCommand(
		"worker",
		"Install CryoSPARC worker",
		runCmd((*command.Command).RunInstallWorker),
	)

	installWorkerCmd.PreRunE = preRunCmd("worker", true)

	installCmd.AddCommand(
		installCompleteCmd,
		installMasterCmd,
		installWorkerCmd,
	)

	// lanes commands
	lanesCmd := newCommand(
		"lanes",
		"Lanes commands",
		runHelpCmd,
	)

	lanesCreateCmd := newCommand(
		"create",
		"Create a new CryoSPARC lane",
		runCmd((*command.Command).RunLanesCreateCmd),
	)

	lanesCreateCmd.PreRunE = preRunCmd("worker", false)

	lanesCreateDefaultCmd := newCommand(
		"default",
		"Install the default CryoSPARC lanes",
		runCmd((*command.Command).RunLanesCreateDefaultCmd),
	)

	lanesCreateDefaultCmd.PreRunE = preRunCmd("worker", false)

	lanesCreateCmd.AddCommand(
		lanesCreateDefaultCmd,
	)

	lanesInstallCmd := newCommand(
		"install",
		"Install a CryoSPARC lane",
		runCmd((*command.Command).RunLanesInstallCmd),
	)

	lanesRemoveCmd := newCommand(
		"remove",
		"Remove a CryoSPARC lane",
		runCmd((*command.Command).RunLanesRemoveCmd),
	)

	lanesCmd.AddCommand(
		lanesCreateCmd,
		lanesInstallCmd,
		lanesRemoveCmd,
	)

	// user commands
	userCmd := newCommand(
		"user",
		"User commands",
		runHelpCmd,
	)

	userCreateCmd := newCommand(
		"create",
		"Create a new CryoSPARC user",
		runCmd((*command.Command).RunUserCreateCmd),
	)

	userCmd.AddCommand(
		userCreateCmd,
	)

	// Base commands
	cobraCmd.AddCommand(
		cryosparcmCmd,
		cryosparcwCmd,
		installCmd,
		lanesCmd,
		userCmd,
	)

	// =============================================
	// cryosparcm CLI Flags
	// =============================================
	addCryosparcmFlags(cryosparcmCmd)

	// =============================================
	// cryosparcw CLI Flags
	// =============================================
	addCryosparcwFlags(cryosparcwCmd)

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
	// lanes CLI Flags
	// =============================================

	// create
	addCommonLanesFlags(lanesCreateCmd)
	addLanesCreateFlags(lanesCreateCmd)

	// create default
	addCommonLanesFlags(lanesCreateDefaultCmd)
	addLanesCreateDefaultFlags(lanesCreateDefaultCmd)

	// install
	addCommonLanesFlags(lanesInstallCmd)
	addLanesInstallFlags(lanesInstallCmd)

	// remove
	addCommonLanesFlags(lanesRemoveCmd)
	addLanesRemoveFlags(lanesRemoveCmd)

	// =============================================
	// user CLI Flags
	// =============================================

	// create
	addCommonUserFlags(userCreateCmd)
	addUserCreateFlags(userCreateCmd)

	// =============================================
	// Execute app
	// =============================================

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Execute cobra command
	if err := fang.Execute(
		context.Background(),
		cobraCmd,
		fang.WithVersion(version)); err != nil {
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

	// If the user did not change the database path and if the cryosparc path
	// was changed, it will append 'database' to the cryosparc path
	if !cmd.Flags().Changed("dbpath") &&
		cfg.CryosparcPath != defaultCryosparcPath {
		cfg.DbPath = filepath.Join(cfg.CryosparcPath, "database")
	}

}

// =============================================================================
// PreRun Commands
// =============================================================================
func preRunCmd(
	nodeType string,
	install bool,
) func(*cobra.Command, []string) error {

	return func(cmd *cobra.Command, _ []string) error {

		if nodeType == "master" || nodeType == "both" {

			if _, ok := validArch[cfg.ArchMaster]; !ok {
				return fmt.Errorf(
					"invalid value for --arch-master: %q "+
						"(must be one of: 'x86_64', 'aarch64')",
					cfg.ArchMaster,
				)
			}

		}

		if nodeType == "worker" || nodeType == "both" {

			if install {

				if _, ok := validArchWorkerInstall[cfg.ArchWorker]; !ok {
					return fmt.Errorf(
						"invalid value for --arch-worker: %q "+
							"(must be one of: 'x86_64', 'aarch64', 'both')",
						cfg.ArchWorker,
					)
				}

			} else {

				if _, ok := validArch[cfg.ArchWorker]; !ok {
					return fmt.Errorf(
						"invalid value for --arch-worker: %q "+
							"(must be one of: 'x86_64', 'aarch64')",
						cfg.ArchWorker,
					)
				}

			}
		}

		return nil
	}

}

// =============================================================================
// Run Commands
// =============================================================================
func runHelpCmd(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func runCmd(
	fn func(*command.Command) error,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		prepareConfig(cmd)
		inst := command.New(cfg)
		return fn(inst)
	}
}

// For commands that allow the user to fill with arguments that are not
// predefined
func runCmdArgs(
	fn func(*command.Command,
		[]string) error,
) func(*cobra.Command, []string) error {

	return func(cmd *cobra.Command, args []string) error {

		prepareConfig(cmd)
		inst := command.New(cfg)
		return fn(inst, args)

	}

}

// =============================================================================
// Flags helper functions
// =============================================================================
func addCryosparcPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&cfg.CryosparcPath,
		"cryosparc-path",
		"d",
		defaultCryosparcPath,
		"Installation directory for CryoSPARC",
	)
}

func addHostnameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.HostName,
		"hostname",
		defaultHostName,
		"Hostname used to access the CryoSPARC web interface",
	)
}

func addCryosparcmHelpFlag(cmd *cobra.Command, command string) {
	cmd.Flags().BoolVar(
		&cfg.CryosparcmHelp,
		"list-commands",
		false,
		fmt.Sprintf("List %s commands\nSame as running "+
			"'%s --help' (only available from version 5)", command, command),
	)
}

func addArchWorkerFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.ArchWorker,
		"arch-worker",
		defaultArchWorker,
		"Worker architecture ('x86_64', 'aarch64')",
	)
}

func addArchWorkerInstallFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.ArchWorker,
		"arch-worker",
		defaultArchWorker,
		"Worker architecture ('x86_64', 'aarch64', 'both')",
	)
}

func addCryosparcLicenseFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.License,
		"license",
		"",
		"CryoSPARC license ID",
	)
}

func addCryosparcVersionFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Version,
		"version",
		defaultVersion,
		"CryoSPARC version to install",
	)
}

func addSsdPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.SSDPath,
		"ssdpath",
		defaultSSDPath,
		"SSD cache directory",
	)
}

func addBasePortFlag(cmd *cobra.Command) {
	cmd.Flags().UintVar(
		&cfg.BasePort,
		"port",
		0,
		"Base TCP port for CryoSPARC services",
	)
}

func addArchMasterFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.ArchMaster,
		"arch-master",
		defaultArchMaster,
		"Master architecture ('x86_64', 'aarch64')",
	)
}

func addCryosparcDbPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.DbPath,
		"dbpath",
		defaultDbPath,
		"CryoSPARC database directory",
	)
}

func addLanesNameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Name,
		"name",
		"",
		"Lane name",
	)
}

func addLanesCachePathFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.CachePath,
		"cache-path",
		defaultLaneCachePath,
		"Cache directory",
	)
}

func addLanesClusterFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Cluster,
		"cluster",
		"",
		"Cluster name (merlin7 or gmerlin7)",
	)
}

func addLanesMemoryFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Memory,
		"memory",
		defaultLaneMemory,
		"Memory requested for the job (in Gb)",
	)
}

func addLanesTimeFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Time,
		"time",
		"",
		"Time requested for the job (format: dd-hh:mm:ss)",
	)
}

func addLanesPartitionFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Partition,
		"partition",
		"",
		"Partition requested for the job (e.g., cpu-hourly or gpu-hourly)",
	)
}

func addLanesGpusFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Gpus,
		"gpus",
		defaultLaneGpus,
		"Number of GPUs requested for the job",
	)
}

func addLanesCpusPerTaskFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.CpusPerTask,
		"cpus-per-task",
		defaultLaneCpusPerTask,
		"Number of CPUs per task requested for "+
			"the job",
	)
}

func addLanesInfoFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Info,
		"info",
		"",
		"Directory of the cluster_info.json "+
			"file (only available from version 5)",
	)
}

func addLanesScriptFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.Lanes.Script,
		"script",
		"",
		"Directory of the cluster_script.sh "+
			"file (only available from version 5)",
	)
}

func addLanesInstallAllFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&cfg.Lanes.InstallAll,
		"all",
		false,
		"Install all lanes stored in the CryoSPARC directory",
	)
}

func addUserEmailFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.User.Email,
		"email",
		"",
		"New user's email",
	)
}

func addUserUsernameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.User.Username,
		"username",
		"",
		"New user's username",
	)
}

func addUserFirstNameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.User.FirstName,
		"firstname",
		"",
		"New user's first name",
	)
}

func addUserLastNameFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&cfg.User.LastName,
		"lastname",
		"",
		"New user's last name",
	)
}

// cryosparcm
func addCryosparcmFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addCryosparcPathFlag(cmd)
	addHostnameFlag(cmd)
	addCryosparcmHelpFlag(cmd, "cryosparcm")

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

// cryosparcw
func addCryosparcwFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addCryosparcPathFlag(cmd)
	addCryosparcmHelpFlag(cmd, "cryosparcw")
	addArchWorkerFlag(cmd)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

// Install
func addCommonInstallFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addCryosparcPathFlag(cmd)
	addCryosparcLicenseFlag(cmd)
	addCryosparcVersionFlag(cmd)

	_ = cmd.MarkFlagRequired("license")

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addInstallMasterFlags(cmd *cobra.Command) {
	addHostnameFlag(cmd)
	addSsdPathFlag(cmd)
	addBasePortFlag(cmd)
	addArchMasterFlag(cmd)
}

func addInstallWorkerFlags(cmd *cobra.Command) {
	addArchWorkerInstallFlag(cmd)
}

func addInstallCompleteFlags(cmd *cobra.Command) {
	addCryosparcDbPathFlag(cmd)
	addHostnameFlag(cmd)
	addBasePortFlag(cmd)
	addArchMasterFlag(cmd)
	addArchWorkerInstallFlag(cmd)
}

// lanes
func addCommonLanesFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addCryosparcPathFlag(cmd)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addLanesCreateFlags(cmd *cobra.Command) {
	addLanesNameFlag(cmd)
	addLanesCachePathFlag(cmd)
	addLanesClusterFlag(cmd)
	addLanesMemoryFlag(cmd)
	addLanesTimeFlag(cmd)
	addLanesPartitionFlag(cmd)
	addLanesGpusFlag(cmd)
	addLanesCpusPerTaskFlag(cmd)
	addArchWorkerFlag(cmd)

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("partition")
	_ = cmd.MarkFlagRequired("time")
	_ = cmd.MarkFlagRequired("cluster")
	_ = cmd.MarkFlagRequired("arch-worker")
}

func addLanesCreateDefaultFlags(cmd *cobra.Command) {
	addArchWorkerFlag(cmd)

	_ = cmd.MarkFlagRequired("arch-worker")
}

func addLanesInstallFlags(cmd *cobra.Command) {
	addLanesNameFlag(cmd)
	addLanesInfoFlag(cmd)
	addLanesScriptFlag(cmd)
	addLanesInstallAllFlag(cmd)

	// Since the lanes can be install with just --name or (--info + --script)
	// but not both. Also, the user can choose to install all lanes with --all
	// which should not conflict with the other flags.
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		nameSet := cmd.Flags().Changed("name")
		allSet := cmd.Flags().Changed("all")
		infoSet := cmd.Flags().Changed("info")
		scriptSet := cmd.Flags().Changed("script")

		if nameSet && allSet {
			return fmt.Errorf(
				"flags '--name' and '--all' cannot be set simultaneously")
		}

		if (infoSet && !scriptSet) || (!infoSet && scriptSet) {
			return fmt.Errorf(
				"both '--info' and '--script' need to be set")
		}

		if (nameSet || allSet) && (infoSet && scriptSet) {
			return fmt.Errorf(
				"set '--info' and '--script' or just '--name' or '--all'")
		}

		if !nameSet && !allSet && !infoSet && !scriptSet {
			return fmt.Errorf(`no lanes selected.
		Use '--name [lane name]' to install a lane stored
		in the CryoSPARC directory. Use '--all' to install all lanes stored
		in the CryoSPARC directory. Or only provide '--info' and '--script'
		to install a lane using cluster_info.json and cluster_script.sh files`)
		}

		return nil
	}
}

func addLanesRemoveFlags(cmd *cobra.Command) {
	addLanesNameFlag(cmd)

	_ = cmd.MarkFlagRequired("name")
}

// user
func addCommonUserFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	addCryosparcPathFlag(cmd)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addUserCreateFlags(cmd *cobra.Command) {
	addUserEmailFlag(cmd)
	addUserUsernameFlag(cmd)
	addUserFirstNameFlag(cmd)
	addUserLastNameFlag(cmd)
}

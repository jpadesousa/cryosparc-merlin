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

	// Lanes
	defaultLaneMemory      = "ram_gb"
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
				"invalid value for --arch-master: %q "+
					"(must be one of: x86_64, aarch64)",
				cfg.ArchMaster,
			)
		}

		if _, ok := validArch[cfg.ArchWorker]; !ok {
			return fmt.Errorf(
				"invalid value for --arch-worker: %q "+
					"(must be one of: x86_64, aarch64)",
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

	lanesCreateDefaultCmd := newCommand(
		"default",
		"Install the default CryoSPARC lanes",
		runCmd((*command.Command).RunLanesCreateDefaultCmd),
	)

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

	// Base commands
	cobraCmd.AddCommand(
		cryosparcmCmd,
		installCmd,
		lanesCmd,
	)

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
	// lanes CLI Flags
	// =============================================

	// create
	addCommonLanesFlags(lanesCreateCmd)
	addLanesCreateFlags(lanesCreateCmd)

	// create default
	addCommonLanesFlags(lanesCreateDefaultCmd)

	// install
	addCommonLanesFlags(lanesInstallCmd)
	addLanesInstallFlags(lanesInstallCmd)

	// remove
	addCommonLanesFlags(lanesRemoveCmd)
	addLanesRemoveFlags(lanesRemoveCmd)

	// =============================================
	// Execute app
	// =============================================

	// Remove completion and help commands (--help is kept)
	cobraCmd.CompletionOptions.DisableDefaultCmd = true
	cobraCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Execute cobra command
	if err := fang.Execute(context.Background(), cobraCmd); err != nil {
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

// cryosparcm
func addCryosparcmFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVarP(&cfg.CryosparcPath, "cryosparc-path", "d",
		defaultCryosparcPath,
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
	flags.StringVarP(&cfg.CryosparcPath, "cryosparc-path", "d",
		defaultCryosparcPath,
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

// lanes
func addCommonLanesFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVarP(&cfg.CryosparcPath, "cryosparc-path", "d",
		defaultCryosparcPath,
		"Installation directory for CryoSPARC",
	)

	flags.SetInterspersed(false)
	flags.SortFlags = false
}

func addLanesCreateFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.Lanes.Name, "name", "",
		"Lane name",
	)

	flags.StringVar(&cfg.Lanes.CachePath, "cache-path", defaultLaneCachePath,
		"Cache directory",
	)

	flags.StringVar(&cfg.Lanes.Cluster, "cluster", "",
		"Cluster name (merlin7 or gmerlin7)",
	)

	flags.StringVar(&cfg.Lanes.Memory, "memory", defaultLaneMemory,
		"Memory requested for the job (in Gb) (default: set by CryoSPARC)",
	)

	flags.StringVar(&cfg.Lanes.Time, "time", "",
		"Time requested for the job (format: dd-hh:mm:ss)",
	)

	flags.StringVar(&cfg.Lanes.Partition, "partition", "",
		"Partition requested for the job (e.g., cpu-hourly or gpu-hourly)",
	)

	flags.StringVar(&cfg.Lanes.Gpus, "gpus", defaultLaneGpus,
		"Number of GPUs requested for the job (default: set by CryoSPARC)",
	)

	flags.StringVar(&cfg.Lanes.CpusPerTask, "cpus-per-task",
		defaultLaneCpusPerTask,
		"Number of CPUs per task requested for "+
			"the job (default: set by CryoSPARC)",
	)

	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("partition")
	_ = cmd.MarkFlagRequired("time")
	_ = cmd.MarkFlagRequired("cluster")
}

func addLanesInstallFlags(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringVar(&cfg.Lanes.Name, "name", "",
		"Lane name stored in the CryoSPARC directory",
	)

	flags.StringVar(&cfg.Lanes.Info, "info", "",
		"Directory of the cluster_info.json file (only available from version 5)",
	)

	flags.StringVar(&cfg.Lanes.Script, "script", "",
		"Directory of the cluster_script.sh file (only available from version 5)",
	)

	flags.BoolVar(&cfg.Lanes.InstallAll, "all", false,
		"Install all lanes stored in the CryoSPARC directory",
	)

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
	flags := cmd.Flags()

	flags.StringVar(&cfg.Lanes.Name, "name", "",
		"Lane name",
	)

	_ = cmd.MarkFlagRequired("name")
}

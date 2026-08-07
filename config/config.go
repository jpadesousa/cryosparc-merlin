package config

// Flag variables
type Config struct {
	Version        string
	CryosparcPath  string
	DbPath         string
	SSDPath        string
	License        string
	HostName       string
	Arch           string
	ArchMaster     string
	ArchWorker     string
	CryosparcmHelp bool
	BasePort       uint

	LanesCreate struct {
		Name        string
		CachePath   string
		Cluster     string
		Memory      string
		Time        string
		Partition   string
		Gpus        string
		CpusPerTask string
	}
}

package config

// Flag variables
type Config struct {
	Version    string
	HomeDir    string
	DbDir      string
	License    string
	RemoteHost string
	Email      string
	Username   string
	FirstName  string
	LastName   string
	Arch       string
	ArchMaster string
	ArchWorker string
	BasePort   uint

	CryosparcmUpdate struct {
		Version      string
		Check        bool
		List         bool
		Override     bool
		DownloadOnly bool
		SkipDownload bool
	}

	CryosparcmPatch struct {
		Install  bool
		Download bool
		Check    bool
		Force    bool
		Yes      bool
	}
}

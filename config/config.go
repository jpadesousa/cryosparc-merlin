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

	// CryosparcmUser struct {
	// 	Email     string
	// 	Username  string
	// 	FirstName string
	// 	LastName  string
	// }

	// CryosparcmUpdate struct {
	// 	Version      string
	// 	Check        bool
	// 	List         bool
	// 	Override     bool
	// 	DownloadOnly bool
	// 	SkipDownload bool
	// }

	// CryosparcmPatch struct {
	// 	Install  bool
	// 	Download bool
	// 	Check    bool
	// 	Force    bool
	// 	Yes      bool
	// }
}

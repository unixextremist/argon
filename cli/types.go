package cli

type CommandType int

const (
	CommandInstall CommandType = iota
	CommandList
	CommandRemove
	CommandSearch
	CommandHelp
	CommandUpgrade
	CommandUnknown
)

type InstallArgs struct {
	Packages []string
	Local    bool
	Branch   string
	Patches  string
	Yes      bool
	PkgDeps  string
	Static   bool
}

type RemoveArgs struct {
	Package string
}

type SearchArgs struct {
	Query string
}

type UpgradeArgs struct {
	Local bool
	Yes   bool
}

package cli

type CommandType int

const (
	CommandInstall CommandType = iota
	CommandList
	CommandRemove
	CommandSearch
	CommandHelp
	CommandUnknown
)

type InstallArgs struct {
	Packages      []string
	Local         bool
	Branch        string
	Patches       string
	Yes           bool
	PkgDeps       string
}

type RemoveArgs struct {
	Package string
}

type SearchArgs struct {
	Query string
}

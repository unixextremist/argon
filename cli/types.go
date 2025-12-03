package cli
type CommandType int
const (
	CommandInstall CommandType = iota
	CommandUnknown
)
type InstallArgs struct {
	Packages      []string
	GitLab        bool
	Codeberg      bool
	Local         bool
	Branch        string
	Patches       string
	Yes           bool
	PkgDeps       string
}
type CliArgs struct {
	Command     CommandType
	InstallArgs InstallArgs
}

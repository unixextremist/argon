package cli

import (
	"flag"
	"os"
	"strings"
)

type CliArgs struct {
	Command     CommandType
	InstallArgs InstallArgs
	RemoveArgs  RemoveArgs
	SearchArgs  SearchArgs
	UpgradeArgs UpgradeArgs
}

func ParseCLI(args []string) CliArgs {
	var cliArgs CliArgs
	if len(args) == 0 {
		cliArgs.Command = CommandUnknown
		return cliArgs
	}

	switch args[0] {
	case "install":
		cliArgs.Command = CommandInstall
		installCmd := flag.NewFlagSet("install", flag.ExitOnError)
		
		local := installCmd.Bool("local", false, "Install locally (~/.local/bin)")
		branch := installCmd.String("branch", "", "Use specific git branch")
		patches := installCmd.String("patches", "", "Apply patches from directory")
		yes := installCmd.Bool("yes", false, "Skip confirmation prompts")
		pkgdeps := installCmd.String("pkgdeps", "", "Install packages from file")
		static := installCmd.Bool("static", false, "Build static binary")
		
		installCmd.Parse(args[1:])
		
		packages := installCmd.Args()
		
		if *pkgdeps != "" {
			content, err := os.ReadFile(*pkgdeps)
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "--") {
						packages = append(packages, line)
					}
				}
			}
		}
		
		cliArgs.InstallArgs = InstallArgs{
			Packages: packages,
			Local:    *local,
			Branch:   *branch,
			Patches:  *patches,
			Yes:      *yes,
			PkgDeps:  *pkgdeps,
			Static:   *static,
		}

	case "list":
		cliArgs.Command = CommandList
	case "remove":
		cliArgs.Command = CommandRemove
		removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
		removeCmd.Parse(args[1:])
		if len(removeCmd.Args()) > 0 {
			cliArgs.RemoveArgs.Package = removeCmd.Args()[0]
		}
	case "search":
		cliArgs.Command = CommandSearch
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		searchCmd.Parse(args[1:])
		if len(searchCmd.Args()) > 0 {
			cliArgs.SearchArgs.Query = strings.Join(searchCmd.Args(), " ")
		}
	case "help":
		cliArgs.Command = CommandHelp
	case "upgrade":
		cliArgs.Command = CommandUpgrade
		upgradeCmd := flag.NewFlagSet("upgrade", flag.ExitOnError)
		local := upgradeCmd.Bool("local", false, "Upgrade local installations only")
		yes := upgradeCmd.Bool("yes", false, "Skip confirmation prompts")
		upgradeCmd.Parse(args[1:])
		cliArgs.UpgradeArgs = UpgradeArgs{
			Local: *local,
			Yes:   *yes,
		}
	default:
		if args[0] == "--help" || args[0] == "-h" {
			cliArgs.Command = CommandHelp
		} else {
			cliArgs.Command = CommandUnknown
		}
	}
	return cliArgs
}

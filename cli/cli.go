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
		local := installCmd.Bool("local", false, "")
		branch := installCmd.String("branch", "", "")
		patches := installCmd.String("patches", "", "")
		yes := installCmd.Bool("yes", false, "")
		pkgdeps := installCmd.String("pkgdeps", "", "")
		installCmd.Parse(args[1:])
		parsedArgs := InstallArgs{
			Local:    *local,
			Branch:   *branch,
			Patches:  *patches,
			Yes:      *yes,
			PkgDeps:  *pkgdeps,
		}
		if *pkgdeps != "" {
			content, err := os.ReadFile(*pkgdeps)
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						parsedArgs.Packages = append(parsedArgs.Packages, line)
					}
				}
			}
		} else {
			parsedArgs.Packages = installCmd.Args()
		}
		cliArgs.InstallArgs = parsedArgs
	case "list":
		cliArgs.Command = CommandList
	case "remove":
		cliArgs.Command = CommandRemove
		if len(args) < 2 {
			cliArgs.Command = CommandUnknown
			return cliArgs
		}
		cliArgs.RemoveArgs.Package = args[1]
	case "search":
		cliArgs.Command = CommandSearch
		if len(args) < 2 {
			cliArgs.Command = CommandUnknown
			return cliArgs
		}
		cliArgs.SearchArgs.Query = args[1]
	case "help":
		cliArgs.Command = CommandHelp
	default:
		cliArgs.Command = CommandUnknown
	}
	return cliArgs
}

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
		
		var packages []string
		var local bool
		var branch string
		var patches string
		var yes bool
		var pkgdeps string
		var static bool
		
		i := 1
		for i < len(args) {
			arg := args[i]
			
			switch arg {
			case "--local":
				local = true
				i++
			case "-local":
				local = true
				i++
			case "--branch":
				if i+1 < len(args) {
					branch = args[i+1]
					i += 2
				} else {
					i++
				}
			case "-branch":
				if i+1 < len(args) {
					branch = args[i+1]
					i += 2
				} else {
					i++
				}
			case "--patches":
				if i+1 < len(args) {
					patches = args[i+1]
					i += 2
				} else {
					i++
				}
			case "-patches":
				if i+1 < len(args) {
					patches = args[i+1]
					i += 2
				} else {
					i++
				}
			case "--yes":
				yes = true
				i++
			case "-yes":
				yes = true
				i++
			case "--pkgdeps":
				if i+1 < len(args) {
					pkgdeps = args[i+1]
					i += 2
				} else {
					i++
				}
			case "-pkgdeps":
				if i+1 < len(args) {
					pkgdeps = args[i+1]
					i += 2
				} else {
					i++
				}
			case "--static":
				static = true
				i++
			case "-static":
				static = true
				i++
			default:
				if strings.HasPrefix(arg, "--branch=") {
					branch = strings.TrimPrefix(arg, "--branch=")
					i++
				} else if strings.HasPrefix(arg, "-branch=") {
					branch = strings.TrimPrefix(arg, "-branch=")
					i++
				} else if strings.HasPrefix(arg, "--patches=") {
					patches = strings.TrimPrefix(arg, "--patches=")
					i++
				} else if strings.HasPrefix(arg, "-patches=") {
					patches = strings.TrimPrefix(arg, "-patches=")
					i++
				} else if strings.HasPrefix(arg, "--pkgdeps=") {
					pkgdeps = strings.TrimPrefix(arg, "--pkgdeps=")
					i++
				} else if strings.HasPrefix(arg, "-pkgdeps=") {
					pkgdeps = strings.TrimPrefix(arg, "-pkgdeps=")
					i++
				} else if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
					i++
				} else {
					packages = append(packages, arg)
					i++
				}
			}
		}
		
		parsedArgs := InstallArgs{
			Local:    local,
			Branch:   branch,
			Patches:  patches,
			Yes:      yes,
			PkgDeps:  pkgdeps,
			Static:   static,
			Packages: packages,
		}

		if pkgdeps != "" {
			content, err := os.ReadFile(pkgdeps)
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(content)), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "--") {
						parsedArgs.Packages = append(parsedArgs.Packages, line)
					}
				}
			}
		}
		cliArgs.InstallArgs = parsedArgs

	case "list":
		cliArgs.Command = CommandList
	case "remove":
		cliArgs.Command = CommandRemove
		removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
		removeCmd.Parse(args[1:])
		cliArgs.RemoveArgs.Package = removeCmd.Arg(0)
	case "search":
		cliArgs.Command = CommandSearch
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		searchCmd.Parse(args[1:])
		cliArgs.SearchArgs.Query = searchCmd.Arg(0)
	case "help":
		cliArgs.Command = CommandHelp
	case "upgrade":
		cliArgs.Command = CommandUpgrade
		upgradeCmd := flag.NewFlagSet("upgrade", flag.ExitOnError)
		local := upgradeCmd.Bool("local", false, "Upgrade local installations")
		yes := upgradeCmd.Bool("yes", false, "Skip confirmation prompts")
		upgradeCmd.Parse(args[1:])
		cliArgs.UpgradeArgs = UpgradeArgs{
			Local: *local,
			Yes:   *yes,
		}
	default:
		cliArgs.Command = CommandUnknown
	}
	return cliArgs
}

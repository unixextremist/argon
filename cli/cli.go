package cli
import (
	"flag"
	"os"
	"strings"
)
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
		gitlab := installCmd.Bool("gitlab", false, "")
		codeberg := installCmd.Bool("codeberg", false, "")
		local := installCmd.Bool("local", false, "")
		branch := installCmd.String("branch", "", "")
		patches := installCmd.String("patches", "", "")
		yes := installCmd.Bool("yes", false, "")
		pkgdeps := installCmd.String("pkgdeps", "", "")
		installCmd.Parse(args[1:])
		parsedArgs := InstallArgs{
			GitLab:   *gitlab,
			Codeberg: *codeberg,
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
	default:
		cliArgs.Command = CommandUnknown
	}
	return cliArgs
}

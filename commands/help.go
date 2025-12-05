package commands

import (
	"fmt"
	"os"
)

func Help() {
	fmt.Println("Argon Package Manager")
	fmt.Println()
	fmt.Println("Usage: argon <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install <package> [options]  Install a package (requires sudo)")
	fmt.Println("  list                          List installed packages")
	fmt.Println("  remove <package>              Remove a package (requires sudo)")
	fmt.Println("  search <query>                Search for packages")
	fmt.Println("  upgrade                       Upgrade installed packages (requires sudo)")
	fmt.Println("  help                          Display this help message")
	fmt.Println()
	fmt.Println("For help with a specific command:")
	fmt.Println("  argon install --help")
	fmt.Println("  argon remove --help")
	fmt.Println("  argon upgrade --help")
	fmt.Println("  argon search --help")
	
	if len(os.Args) > 2 {
		cmd := os.Args[1]
		arg := os.Args[2]
		if arg == "--help" {
			switch cmd {
			case "install":
				fmt.Println()
				fmt.Println("Install options:")
				fmt.Println("  --local         Install locally (~/.local/bin)")
				fmt.Println("  --branch <br>   Use specific git branch")
				fmt.Println("  --patches <dir> Apply patches from directory")
				fmt.Println("  --yes           Skip confirmation prompts")
				fmt.Println("  --pkgdeps <file> Install packages from file")
				fmt.Println("  --static        Build static binary")
			case "upgrade":
				fmt.Println()
				fmt.Println("Upgrade options:")
				fmt.Println("  --local         Upgrade local installations")
				fmt.Println("  --yes           Skip confirmation prompts")
			case "remove":
				fmt.Println()
				fmt.Println("Remove options:")
				fmt.Println("  <package>       Package name to remove")
			case "search":
				fmt.Println()
				fmt.Println("Search options:")
				fmt.Println("  <query>         Search query")
			}
		}
	}
}

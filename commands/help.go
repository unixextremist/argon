package commands

import (
	"fmt"
)

func Help() {
	fmt.Println("Argon Package Manager")
	fmt.Println("Usage: sudo argon <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  install <package> [options]  Install a package")
	fmt.Println("  list                          List installed packages")
	fmt.Println("  remove <package>              Remove a package")
	fmt.Println("  help                          Display this help message")
	fmt.Println()
	fmt.Println("Install options:")
	fmt.Println("  --local         Install locally (~/.local/bin)")
	fmt.Println("  --branch <br>   Use specific git branch")
	fmt.Println("  --patches <dir> Apply patches from directory")
	fmt.Println("  --yes           Skip confirmation prompts")
	fmt.Println("  --pkgdeps <file> Install packages from file")
}

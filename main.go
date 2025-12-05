package main

import (
	"fmt"
	"os"
	"os/user"
	"argon-go/cli"
	"argon-go/commands"
	"argon-go/utils"
)

func main() {
	utils.SetupArgonDirs()
	args := cli.ParseCLI(os.Args[1:])
	if args.InstallArgs.Static {
		fmt.Println("DEBUG CLI Static=true")
	}
	switch args.Command {
	case cli.CommandInstall:
		currentUser, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
			os.Exit(1)
		}
		if currentUser.Uid != "0" {
			fmt.Fprintf(os.Stderr, "Error: argon install must be run with sudo\n")
			fmt.Fprintf(os.Stderr, "Usage: sudo argon install <package> [options]\n")
			os.Exit(1)
		}
		commands.HandleInstall(&args.InstallArgs)
	case cli.CommandList:
		commands.List()
	case cli.CommandRemove:
		currentUser, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
			os.Exit(1)
		}
		if currentUser.Uid != "0" {
			fmt.Fprintf(os.Stderr, "Error: argon remove must be run with sudo\n")
			fmt.Fprintf(os.Stderr, "Usage: sudo argon remove <package>\n")
			os.Exit(1)
		}
		commands.Remove(args.RemoveArgs.Package)
	case cli.CommandSearch:
		commands.Search(args.SearchArgs.Query)
	case cli.CommandHelp:
		commands.Help()
	case cli.CommandUpgrade:
		currentUser, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
			os.Exit(1)
		}
		if currentUser.Uid != "0" {
			fmt.Fprintf(os.Stderr, "Error: argon upgrade must be run with sudo\n")
			fmt.Fprintf(os.Stderr, "Usage: sudo argon upgrade [options]\n")
			os.Exit(1)
		}
		commands.HandleUpgrade(&args.UpgradeArgs)
	default:
		fmt.Println("Usage: argon <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  install <package> [options]  Install a package (requires sudo)")
		fmt.Println("  list                          List installed packages")
		fmt.Println("  remove <package>              Remove a package (requires sudo)")
		fmt.Println("  search <query>                Search for packages")
		fmt.Println("  upgrade                       Upgrade installed packages (requires sudo)")
		fmt.Println("  help                          Display this help message")
		os.Exit(1)
	}
}

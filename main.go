package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"argon-go/cli"
	"argon-go/commands"
	"argon-go/utils"
)

func requireRoot() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current user: %v\n", err)
		os.Exit(1)
	}
	if currentUser.Uid != "0" {
		fmt.Fprintf(os.Stderr, "Error: this command must be run with sudo\n")
		os.Exit(1)
	}
}

func main() {
	utils.SetupArgonDirs()
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nInterrupted. Cleaning up...")
		cancel()
		os.Exit(1)
	}()

	args := cli.ParseCLI(os.Args[1:])
	
	switch args.Command {
	case cli.CommandInstall:
		requireRoot()
		commands.HandleInstall(ctx, &args.InstallArgs)
	case cli.CommandList:
		commands.List()
	case cli.CommandRemove:
		requireRoot()
		commands.Remove(args.RemoveArgs.Package)
	case cli.CommandSearch:
		commands.Search(args.SearchArgs.Query)
	case cli.CommandHelp:
		commands.Help(os.Args)
	case cli.CommandUpgrade:
		requireRoot()
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

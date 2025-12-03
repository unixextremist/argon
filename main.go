package main
import (
	"fmt"
	"os"
	"argon-go/cli"
	"argon-go/commands"
	"argon-go/utils"
)
func main() {
	utils.SetupArgonDirs()
	args := cli.ParseCLI(os.Args[1:])
	if args.Command == cli.CommandInstall {
		commands.HandleInstall(&args.InstallArgs)
	} else {
		fmt.Println("Usage: argon install <package> [options]")
		os.Exit(1)
	}
}

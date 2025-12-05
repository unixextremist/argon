package commands

import (
	"fmt"
	"argon-go/utils"
)

func List() {
	fmt.Println("Installed packages:")

	packages := utils.GetInstalledPackages()
	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return
	}

	for _, pkg := range packages {
		fmt.Printf("- %s\n", pkg.Name)
		fmt.Printf("  Repo: %s\n", pkg.Repo)
		fmt.Printf("  Build System: %s\n", pkg.BuildSystem)
		fmt.Printf("  Hash: %s\n", pkg.Hash)
	}
}

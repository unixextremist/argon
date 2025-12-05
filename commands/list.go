package commands

import (
	"fmt"
	"argon-go/utils"
)

func List() {
	packages := utils.GetInstalledPackages()
	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return
	}
	fmt.Println("Installed packages:")
	for _, pkg := range packages {
		loc := "system"
		if pkg.Local {
			loc = "local"
		}
		fmt.Printf("  %-25s  %-15s  %-6s  %s\n", pkg.Name, pkg.BuildSystem, loc, pkg.Hash[:8])
	}
}

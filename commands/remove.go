package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"argon-go/utils"
)

func Remove(packageName string) {
	installed := utils.GetInstalledPackages()

	for i, pkg := range installed {
		if pkg.Name == packageName {
			installed = append(installed[:i], installed[i+1:]...)
			
			if err := utils.SaveInstalledPackages(installed); err != nil {
				panic(fmt.Sprintf("Failed to update package list: %v", err))
			}
			
			systemPath := filepath.Join("/usr/local/bin", packageName)
			
			if _, err := os.Stat(systemPath); err == nil {
				if err := os.Remove(systemPath); err != nil {
					fmt.Printf("Warning: Failed to remove system binary: %v\n", err)
				} else {
					fmt.Println("Removed system binary:", systemPath)
				}
			}
			
			fmt.Println("Removed successfully")
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Error: Package '%s' not found\n", packageName)
}

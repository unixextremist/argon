package commands

import (
	"fmt"
	"os"
	"os/exec"
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
			
			localPath := filepath.Join(os.Getenv("HOME"), ".local", "bin", packageName)
			systemPath := filepath.Join("/usr/local/bin", packageName)
			
			if _, err := os.Stat(localPath); err == nil {
				if err := os.Remove(localPath); err != nil {
					fmt.Printf("Warning: Failed to remove local binary: %v\n", err)
				} else {
					fmt.Println("Removed local binary:", localPath)
				}
			}
			
			if _, err := os.Stat(systemPath); err == nil {
				cmd := exec.Command("rm", "-f", systemPath)
				if err := cmd.Run(); err != nil {
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

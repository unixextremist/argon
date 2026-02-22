package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"argon-go/utils"
)

func Remove(packageName string) {
	if packageName == "" {
		fmt.Println("Error: no package specified")
		return
	}

	packages := utils.GetInstalledPackages()
	
	found := false
	var updatedPackages []utils.Package
	var pkgToRemove utils.Package
	
	for _, pkg := range packages {
		if pkg.Name == packageName {
			found = true
			pkgToRemove = pkg
		} else {
			updatedPackages = append(updatedPackages, pkg)
		}
	}
	
	if !found {
		fmt.Printf("Package %s not found\n", packageName)
		return
	}
	
	var destPath string
	if pkgToRemove.Local {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		destPath = filepath.Join(home, ".local", "bin", pkgToRemove.Name)
	} else {
		destPath = filepath.Join("/usr/local/bin", pkgToRemove.Name)
	}
	
	if _, err := os.Stat(destPath); err == nil {
		if err := os.Remove(destPath); err != nil {
			fmt.Printf("Error removing binary: %v\n", err)
			return
		}
		fmt.Printf("Removed binary: %s\n", destPath)
	} else {
		fmt.Printf("Binary not found: %s\n", destPath)
	}
	
	buildDir := filepath.Join("/tmp/argon/builds", pkgToRemove.Name)
	if utils.DirectoryExists(buildDir) {
		if err := os.RemoveAll(buildDir); err != nil {
			fmt.Printf("Warning: could not remove build directory: %v\n", err)
		} else {
			fmt.Printf("Removed build directory: %s\n", buildDir)
		}
	}
	
	if err := utils.SaveInstalledPackages(updatedPackages); err != nil {
		fmt.Printf("Error updating package list: %v\n", err)
		return
	}
	
	fmt.Printf("Removed %s\n", packageName)
}

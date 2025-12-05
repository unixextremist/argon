package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"argon-go/utils"
)

func Remove(packageName string) {
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
		home, _ := os.UserHomeDir()
		destPath = filepath.Join(home, ".local", "bin", pkgToRemove.Name)
	} else {
		destPath = filepath.Join("/usr/local/bin", pkgToRemove.Name)
	}
	
	if _, err := os.Stat(destPath); err == nil {
		cmd := exec.Command("rm", "-f", destPath)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error removing binary: %v\n", err)
			return
		}
		fmt.Printf("Removed binary: %s\n", destPath)
	}
	
	buildDir := filepath.Join("/tmp/argon/builds", pkgToRemove.Name)
	if _, err := os.Stat(buildDir); err == nil {
		os.RemoveAll(buildDir)
		fmt.Printf("Removed build directory: %s\n", buildDir)
	}
	
	if err := utils.SaveInstalledPackages(updatedPackages); err != nil {
		fmt.Printf("Error updating package list: %v\n", err)
		return
	}
	
	fmt.Printf("Removed %s\n", packageName)
}

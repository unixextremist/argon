package commands

import (
	"context"
	"fmt"
	"argon-go/cli"
	"argon-go/utils"
)

func checkForUpdate(pkg utils.Package) (bool, string, error) {
	currentHash := pkg.Hash
	remoteHash, err := utils.GetRemoteHash(pkg.Repo, "")
	if err != nil {
		return false, "", err
	}
	return remoteHash != currentHash, remoteHash, nil
}

func upgradePackage(pkg utils.Package, yes bool) {
	hasUpdate, newHash, err := checkForUpdate(pkg)
	if err != nil {
		fmt.Printf("Error checking updates for %s: %v\n", pkg.Name, err)
		return
	}
	if !hasUpdate {
		fmt.Printf("%s is already up to date\n", pkg.Name)
		return
	}
	
	oldHash := pkg.Hash
	if len(oldHash) > 8 {
		oldHash = oldHash[:8]
	}
	newHashShort := newHash
	if len(newHashShort) > 8 {
		newHashShort = newHashShort[:8]
	}
	
	fmt.Printf("Updating %s (%s -> %s)\n", pkg.Name, oldHash, newHashShort)
	
	installArgs := &cli.InstallArgs{
		Packages: []string{pkg.Repo},
		Local:    pkg.Local,
		Yes:      yes,
		Static:   pkg.Static,
	}
	
	ctx := context.Background()
	if err := installSingle(ctx, pkg.Repo, installArgs); err != nil {
		fmt.Printf("Failed to upgrade %s: %v\n", pkg.Name, err)
	}
}

func HandleUpgrade(args *cli.UpgradeArgs) {
	packages := utils.GetInstalledPackages()
	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return
	}
	
	toUpgrade := packages
	
	if args.Local {
		var filtered []utils.Package
		for _, pkg := range packages {
			if pkg.Local {
				filtered = append(filtered, pkg)
			}
		}
		toUpgrade = filtered
	}
	
	if len(toUpgrade) == 0 {
		if args.Local {
			fmt.Println("No local packages to upgrade")
		} else {
			fmt.Println("No packages to upgrade")
		}
		return
	}
	
	fmt.Printf("Found %d packages to upgrade\n", len(toUpgrade))
	for i, pkg := range toUpgrade {
		fmt.Printf("\n[%d/%d] ", i+1, len(toUpgrade))
		upgradePackage(pkg, args.Yes)
	}
	fmt.Println("\nUpgrade complete")
}

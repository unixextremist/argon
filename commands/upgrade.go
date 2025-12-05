package commands

import (
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
	fmt.Printf("Updating %s (%s -> %s)\n", pkg.Name, pkg.Hash[:8], newHash[:8])
	installArgs := &cli.InstallArgs{
		Packages: []string{pkg.Repo},
		Local:    pkg.Local,
		Yes:      yes,
		Static:   false,
	}
	installSingle(pkg.Repo, installArgs)
}

func HandleUpgrade(args *cli.UpgradeArgs) {
	packages := utils.GetInstalledPackages()
	if len(packages) == 0 {
		fmt.Println("No packages installed")
		return
	}
	toUpgrade := packages
	if len(toUpgrade) == 0 {
		fmt.Println("No packages to upgrade")
		return
	}
	fmt.Printf("Found %d packages to upgrade\n", len(toUpgrade))
	for i, pkg := range toUpgrade {
		fmt.Printf("\n[%d/%d] ", i+1, len(toUpgrade))
		upgradePackage(pkg, args.Yes)
	}
	fmt.Println("\nUpgrade complete")
}

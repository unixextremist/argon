package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"argon-go/cli"
	"argon-go/pkgconfig"
	"argon-go/utils"
)

func runCommand(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func handleExistingDir(buildDir string) bool {
	if !utils.DirectoryExists(buildDir) {
		return false
	}
	if utils.IsDirEmpty(buildDir) {
		return false
	}
	fmt.Printf("\nBuild directory '%s' already exists.\n", buildDir)
	fmt.Println("Choose an option:")
	fmt.Println("  1. Use existing directory")
	fmt.Println("  2. Remove directory and re-clone")
	fmt.Println("  3. Abort installation")
	fmt.Print("Choice [1-3]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	switch choice {
	case "1":
		fmt.Println("Using existing directory...")
		return true
	case "2":
		fmt.Println("Removing directory and re-cloning...")
		os.RemoveAll(buildDir)
		return false
	default:
		fmt.Println("Installation aborted")
		return false
	}
}

func cloneRepo(pkg, branch, buildDir string) error {
	domain := utils.GetDomainFromURL(pkg)
	repoPath := utils.ExtractRepoPath(pkg)
	var cmd string
	if branch != "" {
		cmd = fmt.Sprintf("git clone --depth=1 --branch %s https://%s/%s %s", branch, domain, repoPath, buildDir)
	} else {
		cmd = fmt.Sprintf("git clone --depth=1 https://%s/%s %s", domain, repoPath, buildDir)
	}
	return runCommand(cmd)
}

func applyPatches(buildDir, patchesDir string) error {
	if patchesDir == "" {
		return nil
	}
	if !utils.DirectoryExists(patchesDir) {
		return fmt.Errorf("patches directory does not exist")
	}
	cmd := fmt.Sprintf("cd %s && find %s -name '*.patch' -exec patch -Np1 -i {} \\;", buildDir, patchesDir)
	return runCommand(cmd)
}

func buildWithMake(buildDir, repoName, cflags, libs string, static bool) (string, error) {
	env := os.Environ()
	if cflags != "" {
		env = append(env, "CFLAGS="+cflags)
	}
	if libs != "" {
		env = append(env, "LDFLAGS="+libs)
	}
	if static {
		staticFlags := "-static"
		if libs != "" {
			staticFlags = libs + " -static"
		}
		env = append(env, "LDFLAGS="+staticFlags)
		env = append(env, "CC=gcc -static")
	}
	cmd := exec.Command("make")
	cmd.Dir = buildDir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "make", cmd.Run()
}

func buildWithCargo(buildDir string, static bool) (string, error) {
	if static {
		return "", fmt.Errorf("static builds not supported for cargo")
	}
	args := []string{"build", "--release"}
	cmd := exec.Command("cargo", args...)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "cargo", cmd.Run()
}

func buildWithCMake(buildDir string, static bool) (string, error) {
	if static {
		return "", fmt.Errorf("static builds not supported for cmake")
	}
	buildPath := filepath.Join(buildDir, "build")
	os.MkdirAll(buildPath, 0755)
	args := []string{".."}
	cmd1 := exec.Command("cmake", args...)
	cmd1.Dir = buildPath
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return "cmake", err
	}
	cmd2 := exec.Command("make")
	cmd2.Dir = buildPath
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	return "cmake", cmd2.Run()
}

func buildWithConfigure(buildDir string, static bool) (string, error) {
	if static {
		return "", fmt.Errorf("static builds not supported for configure")
	}
	args := []string{"./configure"}
	cmd1 := exec.Command(args[0], args[1:]...)
	cmd1.Dir = buildDir
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return "configure", err
	}
	cmd2 := exec.Command("make")
	cmd2.Dir = buildDir
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	return "configure", cmd2.Run()
}

func buildWithZig(buildDir string, static bool) (string, error) {
	if static {
		return "", fmt.Errorf("static builds not supported for zig")
	}
	args := []string{"build"}
	cmd := exec.Command("zig", args...)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "zig", cmd.Run()
}

func buildWithShellScript(buildDir string, static bool) (string, error) {
	if static {
		return "", fmt.Errorf("static builds not supported for shell scripts")
	}
	scriptPath := filepath.Join(buildDir, "build.sh")
	os.Chmod(scriptPath, 0755)
	env := os.Environ()
	cmd := exec.Command("./build.sh")
	cmd.Dir = buildDir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "shell", cmd.Run()
}

func findBuildFilesRecursive(startDir string) ([]string, string) {
	buildFiles := []string{}
	knownFiles := []string{"Makefile", "makefile", "Cargo.toml", "CMakeLists.txt", "configure", "build.zig", "build.sh"}

	currentDir := startDir
	for {
		for _, file := range knownFiles {
			path := filepath.Join(currentDir, file)
			if utils.FileExists(path) {
				buildFiles = append(buildFiles, path)
			}
		}
		if len(buildFiles) > 0 {
			return buildFiles, currentDir
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			break
		}
		currentDir = parent
	}
	return buildFiles, ""
}

func displayBuildFileWithLess(filepath string) error {
	cmd := exec.Command("less", filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func confirmBuild() bool {
	fmt.Print("\nProceed with build? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func detectAndBuild(buildDir, repoName string, static bool) (string, error) {
	fmt.Printf("DEBUG: static=%t\n", static)
	
	if static {
		return "", fmt.Errorf("static builds not supported")
	}

	if !pkgconfig.CheckPkgConfigExists() {
		fmt.Println("Warning: pkg-config not found in PATH")
	}

	cflags, libs := pkgconfig.GetFlags(repoName, static)

	buildFiles, foundDir := findBuildFilesRecursive(buildDir)
	if len(buildFiles) == 0 {
		return "", fmt.Errorf("no supported build system found")
	}

	var selectedBuildFile string
	if len(buildFiles) == 1 {
		selectedBuildFile = buildFiles[0]
	} else {
		fmt.Println("Multiple build files found:")
		for i, file := range buildFiles {
			rel, _ := filepath.Rel(buildDir, file)
			fmt.Printf("%d. %s\n", i+1, rel)
		}
		fmt.Print("Select build file [1]: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		index := 0
		if choice != "" {
			fmt.Sscanf(choice, "%d", &index)
			index--
		}
		if index < 0 || index >= len(buildFiles) {
			index = 0
		}
		selectedBuildFile = buildFiles[index]
	}

	fmt.Printf("Using build file: %s\n", selectedBuildFile)
	fmt.Println("Displaying build file with less (press q to continue)...")
	if err := displayBuildFileWithLess(selectedBuildFile); err != nil {
		fmt.Printf("Warning: could not display with less: %v\n", err)
	}

	if !confirmBuild() {
		return "", fmt.Errorf("build cancelled by user")
	}

	buildDir = foundDir
	filename := filepath.Base(selectedBuildFile)

	switch filename {
	case "Makefile", "makefile":
		return buildWithMake(buildDir, repoName, cflags, libs, static)
	case "Cargo.toml":
		return buildWithCargo(buildDir, static)
	case "CMakeLists.txt":
		return buildWithCMake(buildDir, static)
	case "configure":
		return buildWithConfigure(buildDir, static)
	case "build.zig":
		return buildWithZig(buildDir, static)
	case "build.sh":
		return buildWithShellScript(buildDir, static)
	default:
		return "", fmt.Errorf("unsupported build file: %s", filename)
	}
}

func findBinary(buildDir, repoName string, static bool) (string, error) {
	possiblePaths := []string{
		filepath.Join(buildDir, repoName),
		filepath.Join(buildDir, "target", "release", repoName),
		filepath.Join(buildDir, "target", "x86_64-unknown-linux-musl", "release", repoName),
		filepath.Join(buildDir, "build", repoName),
	}
	
	for _, path := range possiblePaths {
		if utils.FileExists(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("binary not found")
}

func installBinary(buildDir, repoName string, local, yes, static bool) error {
	binaryPath, err := findBinary(buildDir, repoName, static)
	if err != nil {
		return err
	}
	var destPath string
	var cmd string
	if local {
		home, _ := os.UserHomeDir()
		destPath = filepath.Join(home, ".local", "bin", repoName)
		os.MkdirAll(filepath.Dir(destPath), 0755)
		cmd = fmt.Sprintf("install -Dm755 %s %s", binaryPath, destPath)
	} else {
		destPath = filepath.Join("/usr/local/bin", repoName)
		cmd = fmt.Sprintf("install -m755 %s %s", binaryPath, destPath)
	}
	return runCommand(cmd)
}

func addToPackageList(pkg, repoName, buildSystem, hash string, local, static bool) error {
	packages := utils.GetInstalledPackages()

	for i, existingPkg := range packages {
		if existingPkg.Name == repoName {
			packages[i].Repo = pkg
			packages[i].BuildSystem = buildSystem
			packages[i].Hash = hash
			packages[i].Local = local
			packages[i].Static = static
			return utils.SaveInstalledPackages(packages)
		}
	}

	newPackage := utils.Package{
		Name:        repoName,
		Repo:        pkg,
		BuildSystem: buildSystem,
		Hash:        hash,
		Local:       local,
		Static:      static,
	}

	packages = append(packages, newPackage)
	return utils.SaveInstalledPackages(packages)
}

func installSingle(pkg string, args *cli.InstallArgs) {
	if strings.HasPrefix(pkg, "--") {
		fmt.Printf("Warning: Skipping '%s' as it appears to be a flag, not a package\n", pkg)
		return
	}

	fmt.Printf("Installing %s\n", pkg)
	start := time.Now()
	repoName := utils.GetRepoName(pkg)
	buildDir := filepath.Join("/tmp/argon/builds", repoName)

	var hash string
	var err error

	if utils.DirectoryExists(buildDir) && !utils.IsDirEmpty(buildDir) {
		if !handleExistingDir(buildDir) {
			os.RemoveAll(buildDir)
		} else {
			hash, err = utils.GetGitHash(buildDir)
			if err != nil {
				fmt.Printf("Warning: Could not get git hash: %v\n", err)
			}
			goto build
		}
	}

	if err := cloneRepo(pkg, args.Branch, buildDir); err != nil {
		fmt.Printf("Failed to clone: %v\n", err)
		return
	}

	hash, err = utils.GetGitHash(buildDir)
	if err != nil {
		fmt.Printf("Warning: Could not get git hash: %v\n", err)
	}

	if err := applyPatches(buildDir, args.Patches); err != nil {
		fmt.Printf("Failed to apply patches: %v\n", err)
	}

build:
	buildSystem, err := detectAndBuild(buildDir, repoName, args.Static)
	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
		return
	}

	if err := installBinary(buildDir, repoName, args.Local, args.Yes, args.Static); err != nil {
		fmt.Printf("Installation failed: %v\n", err)
		return
	}

	if err := addToPackageList(pkg, repoName, buildSystem, hash, args.Local, args.Static); err != nil {
		fmt.Printf("Warning: Could not update package list: %v\n", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Installed in %.2fs\n", elapsed.Seconds())
}

func HandleInstall(args *cli.InstallArgs) {
	if len(args.Packages) == 0 {
		fmt.Println("No packages specified")
		return
	}

	var actualPackages []string
	for _, pkg := range args.Packages {
		if !strings.HasPrefix(pkg, "--") {
			actualPackages = append(actualPackages, pkg)
		}
	}

	if len(actualPackages) == 0 {
		fmt.Println("No valid packages specified")
		return
	}

	for _, pkg := range actualPackages {
		installSingle(pkg, args)
	}
}

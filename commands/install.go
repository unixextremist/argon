package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"argon-go/cli"
	"argon-go/pkgconfig"
	"argon-go/utils"
)

func cloneRepo(pkg, branch, buildDir string) error {
	domain := utils.GetDomainFromURL(pkg)
	repoPath := utils.ExtractRepoPath(pkg)
	url := fmt.Sprintf("https://%s/%s", domain, repoPath)
	
	args := []string{"clone", "--depth=1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, url, buildDir)
	
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func handleExistingDir(buildDir string) (bool, error) {
	if !utils.DirectoryExists(buildDir) {
		return false, nil
	}
	if utils.IsDirEmpty(buildDir) {
		return false, nil
	}
	
	fmt.Printf("\nBuild directory '%s' already exists.\n", buildDir)
	fmt.Println("Choose an option:")
	fmt.Println("  1. Use existing directory")
	fmt.Println("  2. Remove directory and re-clone")
	fmt.Println("  3. Abort installation")
	fmt.Print("Choice [1-3]: ")
	
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	choice = strings.TrimSpace(choice)
	
	switch choice {
	case "1":
		fmt.Println("Using existing directory...")
		return true, nil
	case "2":
		fmt.Println("Removing directory and re-cloning...")
		if err := os.RemoveAll(buildDir); err != nil {
			return false, fmt.Errorf("failed to remove directory: %w", err)
		}
		return false, nil
	default:
		return false, fmt.Errorf("installation aborted by user")
	}
}

func applyPatches(buildDir, patchesDir string) error {
	if patchesDir == "" {
		return nil
	}
	
	if !utils.DirectoryExists(patchesDir) {
		return fmt.Errorf("patches directory does not exist: %s", patchesDir)
	}
	
	cleanPatchesDir := filepath.Clean(patchesDir)
	if strings.Contains(cleanPatchesDir, "..") {
		return fmt.Errorf("invalid patches directory path")
	}
	
	cmd := exec.Command("find", cleanPatchesDir, "-name", "*.patch", "-type", "f")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find patches: %w", err)
	}
	
	patches := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, patch := range patches {
		if patch == "" {
			continue
		}
		patchCmd := exec.Command("patch", "-Np1", "-i", patch)
		patchCmd.Dir = buildDir
		patchCmd.Stdout = os.Stdout
		patchCmd.Stderr = os.Stderr
		if err := patchCmd.Run(); err != nil {
			return fmt.Errorf("failed to apply patch %s: %w", patch, err)
		}
	}
	
	return nil
}

func buildWithMake(buildDir, repoName, cflags, libs string, static bool) (string, error) {
	env := os.Environ()
	
	if static {
		staticCflags := "-static"
		if cflags != "" {
			staticCflags = cflags + " -static"
		}
		env = append(env, "CFLAGS="+staticCflags)
		
		staticLdflags := "-static"
		if libs != "" {
			staticLdflags = libs + " -static"
		}
		env = append(env, "LDFLAGS="+staticLdflags)
	} else {
		if cflags != "" {
			env = append(env, "CFLAGS="+cflags)
		}
		if libs != "" {
			env = append(env, "LDFLAGS="+libs)
		}
	}
	
	cmd := exec.Command("make")
	cmd.Dir = buildDir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "make", cmd.Run()
}

func buildWithCargo(buildDir string, static bool) (string, error) {
	args := []string{"build", "--release"}
	
	if static {
		target := fmt.Sprintf("%s-unknown-linux-musl", runtime.GOARCH)
		args = append(args, "--target", target)
		
		installTarget := exec.Command("rustup", "target", "add", target)
		installTarget.Stdout = io.Discard
		installTarget.Stderr = io.Discard
		installTarget.Run() 
	}
	
	cmd := exec.Command("cargo", args...)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "cargo", cmd.Run()
}

func buildWithCMake(buildDir string, static bool) (string, error) {
	buildPath := filepath.Join(buildDir, "build")
	if err := os.MkdirAll(buildPath, 0755); err != nil {
		return "cmake", err
	}
	
	cmakeArgs := []string{".."}
	if static {
		cmakeArgs = append(cmakeArgs, "-DCMAKE_EXE_LINKER_FLAGS=-static", "-DBUILD_SHARED_LIBS=OFF")
	}
	
	cmd1 := exec.Command("cmake", cmakeArgs...)
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
	configureArgs := []string{"./configure"}
	if static {
		configureArgs = append(configureArgs, "LDFLAGS=-static", "--disable-shared")
	}
	
	cmd1 := exec.Command(configureArgs[0], configureArgs[1:]...)
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
	args := []string{"build"}
	if static {
		args = append(args, "-Dtarget=native-native-musl")
	}
	
	cmd := exec.Command("zig", args...)
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "zig", cmd.Run()
}

func buildWithShellScript(buildDir string, static bool) (string, error) {
	scriptPath := filepath.Join(buildDir, "build.sh")
	
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return "shell", fmt.Errorf("failed to make build.sh executable: %w", err)
	}
	
	env := os.Environ()
	if static {
		env = append(env, "STATIC_BUILD=1")
	}
	
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
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func detectAndBuild(buildDir, repoName string, static bool) (string, error) {
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
	targetDir := "release"
	if static {
		targetDir = fmt.Sprintf("%s-unknown-linux-musl/release", runtime.GOARCH)
	}
	
	possiblePaths := []string{
		filepath.Join(buildDir, repoName),
		filepath.Join(buildDir, "target", "release", repoName),
		filepath.Join(buildDir, "target", targetDir, repoName),
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
	if local {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		destPath = filepath.Join(home, ".local", "bin", repoName)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else {
		destPath = filepath.Join("/usr/local/bin", repoName)
	}
	
	sourceData, err := os.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to read binary: %w", err)
	}
	
	if err := os.WriteFile(destPath, sourceData, 0755); err != nil {
		return fmt.Errorf("failed to install binary: %w", err)
	}
	
	fmt.Printf("Installed: %s -> %s\n", binaryPath, destPath)
	return nil
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

func installSingle(ctx context.Context, pkg string, args *cli.InstallArgs) error {
	if strings.HasPrefix(pkg, "--") {
		return fmt.Errorf("invalid package name: %s", pkg)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fmt.Printf("Installing %s\n", pkg)
	start := time.Now()
	repoName := utils.GetRepoName(pkg)
	buildDir := filepath.Join("/tmp/argon/builds", repoName)

	var hash string
	var err error

	if utils.DirectoryExists(buildDir) && !utils.IsDirEmpty(buildDir) {
		useExisting, err := handleExistingDir(buildDir)
		if err != nil {
			return err
		}
		if useExisting {
			hash, err = utils.GetGitHash(buildDir)
			if err != nil {
				fmt.Printf("Warning: Could not get git hash: %v\n", err)
			}
			goto build
		}
	}

	if err := cloneRepo(pkg, args.Branch, buildDir); err != nil {
		return fmt.Errorf("failed to clone: %w", err)
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
		return fmt.Errorf("build failed: %w", err)
	}

	if err := installBinary(buildDir, repoName, args.Local, args.Yes, args.Static); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	if err := addToPackageList(pkg, repoName, buildSystem, hash, args.Local, args.Static); err != nil {
		fmt.Printf("Warning: Could not update package list: %v\n", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Installed in %.2fs\n", elapsed.Seconds())
	return nil
}

func HandleInstall(ctx context.Context, args *cli.InstallArgs) {
	if len(args.Packages) == 0 {
		fmt.Println("No packages specified")
		return
	}

	for _, pkg := range args.Packages {
		if err := installSingle(ctx, pkg, args); err != nil {
			fmt.Printf("Error installing %s: %v\n", pkg, err)
		}
	}
}

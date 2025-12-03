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
func getDomain(args *cli.InstallArgs) string {
	if args.Codeberg {
		return "codeberg.org"
	}
	if args.GitLab {
		return "gitlab.com"
	}
	return "github.com"
}
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
func cloneRepo(pkg, domain, branch, buildDir string) error {
	var cmd string
	if branch != "" {
		cmd = fmt.Sprintf("git clone --depth=1 --branch %s https://%s/%s %s", branch, domain, pkg, buildDir)
	} else {
		cmd = fmt.Sprintf("git clone --depth=1 https://%s/%s %s", domain, pkg, buildDir)
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
func buildWithMake(buildDir, repoName, cflags, libs string) error {
	env := os.Environ()
	if cflags != "" {
		env = append(env, "CFLAGS="+cflags)
	}
	if libs != "" {
		env = append(env, "LDFLAGS="+libs)
	}
	cmd := exec.Command("make")
	cmd.Dir = buildDir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func buildWithCargo(buildDir string) error {
	cmd := exec.Command("cargo", "build", "--release")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func buildWithCMake(buildDir string) error {
	buildPath := filepath.Join(buildDir, "build")
	os.MkdirAll(buildPath, 0755)
	cmd1 := exec.Command("cmake", "..")
	cmd1.Dir = buildPath
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return err
	}
	cmd2 := exec.Command("make")
	cmd2.Dir = buildPath
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	return cmd2.Run()
}
func buildWithConfigure(buildDir string) error {
	cmd1 := exec.Command("./configure")
	cmd1.Dir = buildDir
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	if err := cmd1.Run(); err != nil {
		return err
	}
	cmd2 := exec.Command("make")
	cmd2.Dir = buildDir
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	return cmd2.Run()
}
func buildWithZig(buildDir string) error {
	cmd := exec.Command("zig", "build")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func buildWithShellScript(buildDir string) error {
	scriptPath := filepath.Join(buildDir, "build.sh")
	os.Chmod(scriptPath, 0755)
	cmd := exec.Command("./build.sh")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func detectAndBuild(buildDir, repoName string) error {
	cflags, libs := pkgconfig.GetFlags(repoName)
	switch {
	case utils.FileExists(filepath.Join(buildDir, "Makefile")) || utils.FileExists(filepath.Join(buildDir, "makefile")):
		return buildWithMake(buildDir, repoName, cflags, libs)
	case utils.FileExists(filepath.Join(buildDir, "Cargo.toml")):
		return buildWithCargo(buildDir)
	case utils.FileExists(filepath.Join(buildDir, "CMakeLists.txt")):
		return buildWithCMake(buildDir)
	case utils.FileExists(filepath.Join(buildDir, "configure")):
		return buildWithConfigure(buildDir)
	case utils.FileExists(filepath.Join(buildDir, "build.zig")):
		return buildWithZig(buildDir)
	case utils.FileExists(filepath.Join(buildDir, "build.sh")):
		return buildWithShellScript(buildDir)
	default:
		return fmt.Errorf("no supported build system found")
	}
}
func findBinary(buildDir, repoName string) (string, error) {
	possiblePaths := []string{
		filepath.Join(buildDir, repoName),
		filepath.Join(buildDir, "target", "release", repoName),
		filepath.Join(buildDir, "build", repoName),
	}
	for _, path := range possiblePaths {
		if utils.FileExists(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("binary not found")
}
func installBinary(buildDir, repoName string, local, yes bool) error {
	binaryPath, err := findBinary(buildDir, repoName)
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
		priv := utils.GetPrivilegeCommand()
		if priv == "" {
			return fmt.Errorf("neither sudo nor doas found")
		}
		if !yes {
			fmt.Print("Install system-wide? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				return nil
			}
		}
		destPath = filepath.Join("/usr/local/bin", repoName)
		cmd = fmt.Sprintf("%s install -m755 %s %s", priv, binaryPath, destPath)
	}
	return runCommand(cmd)
}
func installSingle(pkg string, args *cli.InstallArgs) {
	fmt.Printf("Installing %s\n", pkg)
	start := time.Now()
	domain := getDomain(args)
	repoName := utils.GetRepoName(pkg)
	buildDir := filepath.Join("/tmp/argon/builds", repoName)
	if utils.DirectoryExists(buildDir) && !utils.IsDirEmpty(buildDir) {
		if !handleExistingDir(buildDir) {
			os.RemoveAll(buildDir)
		} else {
			goto build
		}
	}
	if err := cloneRepo(pkg, domain, args.Branch, buildDir); err != nil {
		fmt.Printf("Failed to clone: %v\n", err)
		return
	}
	if err := applyPatches(buildDir, args.Patches); err != nil {
		fmt.Printf("Failed to apply patches: %v\n", err)
	}
build:
	if err := detectAndBuild(buildDir, repoName); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		return
	}
	if err := installBinary(buildDir, repoName, args.Local, args.Yes); err != nil {
		fmt.Printf("Installation failed: %v\n", err)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("Installed in %.2fs\n", elapsed.Seconds())
}
func HandleInstall(args *cli.InstallArgs) {
	if len(args.Packages) == 0 {
		fmt.Println("No packages specified")
		return
	}
	for i, pkg := range args.Packages {
		installSingle(pkg, args)
		if i < len(args.Packages)-1 {
			fmt.Println()
		}
	}
}

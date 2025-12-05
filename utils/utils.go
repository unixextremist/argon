package utils

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Package struct {
	Name        string `json:"name"`
	Repo        string `json:"repo"`
	BuildSystem string `json:"build_system"`
	Hash        string `json:"hash"`
}

func GetInstalledPackages() []Package {
	filePath := "/var/lib/argon/list"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []Package{}
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		panic(err)
	}
	return packages
}

func SaveInstalledPackages(packages []Package) error {
	filePath := "/var/lib/argon/list"
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func SetupArgonDirs() {
	os.MkdirAll("/tmp/argon/builds", 0755)
	os.MkdirAll("/var/lib/argon", 0755)
}

func GetPrivilegeCommand() string {
	if _, err := exec.LookPath("doas"); err == nil {
		return "doas"
	}
	if _, err := exec.LookPath("sudo"); err == nil {
		return "sudo"
	}
	return ""
}

func GetRepoName(pkg string) string {
	parts := strings.Split(pkg, "/")
	name := parts[len(parts)-1]
	return strings.TrimSuffix(name, ".git")
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirEmpty(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return true
	}
	defer dir.Close()
	entries, _ := dir.Readdir(1)
	return len(entries) == 0
}

func BuildPath(base, part string) string {
	return filepath.Join(base, part)
}

func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

func GetGitHash(buildDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = buildDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

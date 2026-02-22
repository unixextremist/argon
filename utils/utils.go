package utils

import (
	"encoding/json"
	"fmt"
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
	Local       bool   `json:"local"`
	Static      bool   `json:"static"`
}

const (
	ArgonLibDir  = "/var/lib/argon"
	ArgonTempDir = "/tmp/argon"
)

func GetInstalledPackages() []Package {
	filePath := filepath.Join(ArgonLibDir, "list")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []Package{}
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []Package{}
	}
	var packages []Package
	if err := json.Unmarshal(data, &packages); err != nil {
		return []Package{}
	}
	return packages
}

func SaveInstalledPackages(packages []Package) error {
	filePath := filepath.Join(ArgonLibDir, "list")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0600)
}

func SetupArgonDirs() {
	os.MkdirAll(filepath.Join(ArgonTempDir, "builds"), 0755)
	os.MkdirAll(ArgonLibDir, 0755)
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

func GetDomainFromURL(pkg string) string {
	if strings.Contains(pkg, "gitlab.com") {
		return "gitlab.com"
	}
	if strings.Contains(pkg, "codeberg.org") {
		return "codeberg.org"
	}
	if strings.Contains(pkg, "github.com") {
		return "github.com"
	}
	parts := strings.Split(pkg, "/")
	if len(parts) > 2 && strings.Contains(parts[0], ".") {
		return parts[0]
	}
	return "github.com"
}

func ExtractRepoPath(pkg string) string {
	parts := strings.Split(pkg, "://")
	if len(parts) > 1 {
		pkg = parts[1]
	}
	parts = strings.Split(pkg, "/")
	start := 0
	for i, part := range parts {
		if strings.Contains(part, ".") {
			start = i + 1
		}
	}
	return strings.Join(parts[start:], "/")
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func IsDirEmpty(path string) bool {
	dir, err := os.Open(path)
	if err != nil {
		return true
	}
	defer dir.Close()
	entries, err := dir.Readdirnames(1)
	return err != nil || len(entries) == 0
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

func GetRemoteHash(repoURL string, branch string) (string, error) {
	domain := GetDomainFromURL(repoURL)
	repoPath := ExtractRepoPath(repoURL)
	
	args := []string{"ls-remote", fmt.Sprintf("https://%s/%s", domain, repoPath)}
	if branch != "" {
		args = append(args, fmt.Sprintf("refs/heads/%s", branch))
	} else {
		args = append(args, "HEAD")
	}
	
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 && lines[0] != "" {
		parts := strings.Fields(lines[0])
		if len(parts) > 0 {
			return parts[0], nil
		}
	}
	return "", fmt.Errorf("no hash found in remote response")
}

package utils
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)
func SetupArgonDirs() {
	os.MkdirAll("/tmp/argon/builds", 0755)
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

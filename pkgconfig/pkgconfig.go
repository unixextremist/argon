package pkgconfig
import (
	"log"
	"os/exec"
	"strings"
)
func GetFlags(pkg string) (string, string) {
	if !CheckPkgConfigExists() {
		log.Printf("Warning: pkg-config not found in path. Skipping pkg-config flags for %s", pkg)
		return "", ""
	}
	cflagsCmd := exec.Command("pkg-config", "--cflags", pkg)
	cflagsOut, _ := cflagsCmd.Output()
	libsCmd := exec.Command("pkg-config", "--libs", pkg)
	libsOut, _ := libsCmd.Output()
	cflags := strings.TrimSpace(string(cflagsOut))
	libs := strings.TrimSpace(string(libsOut))
	return cflags, libs
}
func CheckPkgConfigExists() bool {
	cmd := exec.Command("pkg-config", "--version")
	return cmd.Run() == nil
}

package pkgconfig

import (
	"os/exec"
	"strings"
)

func CheckPkgConfigExists() bool {
	_, err := exec.LookPath("pkg-config")
	return err == nil
}

func GetFlags(pkg string, static bool) (string, string) {
	args := []string{"--cflags", "--libs"}
	if static {
		args = append([]string{"--static"}, args...)
	}
	
	cmd := exec.Command("pkg-config", append(args, pkg)...)
	output, err := cmd.Output()
	
	if err != nil && static {
		return GetFlags(pkg, false)
	}
	
	if err != nil {
		return "", ""
	}
	
	flags := strings.TrimSpace(string(output))
	if flags == "" {
		return "", ""
	}
	
	return parseFlags(flags, static && !strings.Contains(flags, "-static"))
}

func parseFlags(flags string, addStatic bool) (string, string) {
	var cflags []string
	var libs []string
	parts := strings.Fields(flags)
	
	for _, part := range parts {
		if strings.HasPrefix(part, "-I") || strings.HasPrefix(part, "-D") || strings.HasPrefix(part, "-f") {
			cflags = append(cflags, part)
		} else if strings.HasPrefix(part, "-L") || strings.HasPrefix(part, "-l") || strings.HasPrefix(part, "-Wl,") {
			libs = append(libs, part)
		}
	}
	
	if addStatic {
		libs = append(libs, "-static")
	}
	
	return strings.Join(cflags, " "), strings.Join(libs, " ")
}

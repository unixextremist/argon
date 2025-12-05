package pkgconfig

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckPkgConfigExists() bool {
	_, err := exec.LookPath("pkg-config")
	return err == nil
}

func GetFlags(pkg string, static bool) (string, string) {
	var cmd string
	if static {
		cmd = fmt.Sprintf("pkg-config --static --cflags --libs %s 2>/dev/null", pkg)
		output, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			cmd = fmt.Sprintf("pkg-config --cflags --libs %s", pkg)
			output, err = exec.Command("sh", "-c", cmd).Output()
			if err != nil {
				return "", "-static"
			}
			flags := strings.TrimSpace(string(output))
			if flags == "" {
				return "", "-static"
			}
			return parseFlags(flags, true)
		}
		flags := strings.TrimSpace(string(output))
		if flags == "" {
			return "", "-static"
		}
		return parseFlags(flags, false)
	} else {
		cmd = fmt.Sprintf("pkg-config --cflags --libs %s", pkg)
		output, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			return "", ""
		}
		flags := strings.TrimSpace(string(output))
		if flags == "" {
			return "", ""
		}
		return parseFlags(flags, false)
	}
}

func parseFlags(flags string, addStatic bool) (string, string) {
	var cflags []string
	var libs []string
	parts := strings.Fields(flags)
	
	for _, part := range parts {
		if strings.HasPrefix(part, "-I") || strings.HasPrefix(part, "-D") {
			cflags = append(cflags, part)
		} else if strings.HasPrefix(part, "-L") || strings.HasPrefix(part, "-l") {
			libs = append(libs, part)
		}
	}
	
	if addStatic && !strings.Contains(strings.Join(libs, " "), "-static") {
		libs = append(libs, "-static")
	}
	
	return strings.Join(cflags, " "), strings.Join(libs, " ")
}

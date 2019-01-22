package pgo

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/juju/loggo"
)

var (
	logger         = loggo.GetLogger("pgo")
	ignorePrefixes = []string{"v", "version", "python/", "Python "}
	reIdentifier   = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)$`)
)

// Config .
type Config struct {
	PPrefix      string
	PMirror      string
	VersionsPath string
}

func getConfig() Config {
	cfg := Config{
		PPrefix:      os.Getenv("HOME"),
		PMirror:      "https://www.python.org/ftp/python/",
		VersionsPath: "p/versions",
	}
	if os.Getenv("P_PREFIX") != "" {
		cfg.PPrefix = os.Getenv("P_PREFIX")
	} else {
		logger.Infof("no P_PREFIX defined, using default: %s", cfg.PPrefix)
	}
	if os.Getenv("P_MIRROR") != "" {
		cfg.PMirror = os.Getenv("P_MIRROR")
	} else {
		logger.Debugf("no P_MIRROR defined, using default: %s", cfg.PMirror)
	}
	return cfg
}

func cleanVersionString(vstr string) (string, error) {
	// ignore prefixes and whitespace
	for _, prefix := range ignorePrefixes {
		vstr = strings.TrimPrefix(vstr, prefix)
	}
	vstr = strings.TrimSpace(vstr)

	// must be X.Y.Z format
	fmt.Println("vstr:", vstr)
	if !reIdentifier.MatchString(vstr) {
		return "", fmt.Errorf("version string must be in X.Y.Z format")
	}

	return vstr, nil
}

// GetCurrentVersion returns the currently active python version
func GetCurrentVersion() (string, error) {
	out, err := exec.Command("python", "--version").Output()
	if err != nil {
		return "", fmt.Errorf("no python found")
	}
	return cleanVersionString(string(out))
}

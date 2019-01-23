package pgo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/juju/loggo"
)

var (
	logger              = loggo.GetLogger("pgo")
	ignorePrefixes      = []string{"v", "version", "python/", "Python "}
	reIdentifier        = regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)`)
	errNotInstalled     = fmt.Errorf("version not installed")
	errAlreadyInstalled = fmt.Errorf("version is already installed")
)

const (
	// name of the executable to be managed
	excName = "python"
	// path after prefix where versions whill be stored
	versionsPath = "p/versions/python"
	// minimum allowed version
	minLegalVersion = "2.7.0"
)

// Config .
type Config struct {
	PPrefix      string
	PMirror      string
	VersionsPath string
}

func getConfig() Config {
	cfg := Config{
		PPrefix: os.Getenv("HOME"),
		PMirror: "https://www.python.org/ftp/python/",
	}
	if os.Getenv("P_PREFIX") != "" {
		cfg.PPrefix = os.Getenv("P_PREFIX")
		logger.Debugf("P_PREFIX: %s", cfg.PPrefix)
	} else {
		logger.Infof("no P_PREFIX defined, using default: %s", cfg.PPrefix)
	}
	if os.Getenv("P_MIRROR") != "" {
		cfg.PMirror = os.Getenv("P_MIRROR")
		logger.Debugf("P_MIRROR: %s", cfg.PMirror)
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
	if !reIdentifier.MatchString(vstr) {
		return "", fmt.Errorf("version string must be in X.Y.Z format")
	}

	return vstr, nil
}

func isVersionInstalled(versionStr string) (bool, error) {
	installedVersions, err := GetInstalledVersions()
	if err != nil {
		return false, err
	}
	return stringContains(installedVersions, versionStr), nil
}

// GetCurrentVersion returns the currently active python version
func GetCurrentVersion() (string, error) {
	out, err := exec.Command(excName, "--version").Output()
	if err != nil {
		return "", fmt.Errorf("no python found")
	}
	return cleanVersionString(string(out))
}

// GetAvailableVersions returns the array of available python versions (from the mirror)
func GetAvailableVersions() ([]string, error) {
	cfg := getConfig()

	minVersion, err := semver.Make(minLegalVersion)
	if err != nil {
		return nil, err
	}

	versions, err := getPythonVersions(cfg.PMirror)
	if err != nil {
		return nil, err
	}

	versionStrs := make([]string, 0, len(versions))
	for _, semver := range versions {
		if semver.Compare(minVersion) < 0 {
			continue
		}
		versionStrs = append(versionStrs, semver.String())
	}
	return versionStrs, nil
}

// GetInstalledVersions returns the array of installed python versions
func GetInstalledVersions() ([]string, error) {
	cfg := getConfig()
	versionsDir := filepath.Join(cfg.PPrefix, versionsPath)
	logger.Debugf("versionsDir: %s", versionsDir)
	versions, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	availableVersions := []string{}
	for _, vDir := range versions {
		if !vDir.IsDir() {
			continue
		}
		logger.Debugf("raw version dir: %s", vDir.Name())
		vStr, err := cleanVersionString(vDir.Name())
		if err == nil {
			logger.Debugf("clean version dir: %s", vStr)
			availableVersions = append(availableVersions, vStr)
		}
	}
	return availableVersions, nil
}

// GetLatestVersion returns the latest available version of python
func GetLatestVersion() (string, error) {
	versions, err := GetAvailableVersions()
	if err != nil {
		return "", err
	}
	return versions[len(versions)-1], nil
}

// GetStableVersion returns the latest available stable version of python
func GetStableVersion() (string, error) {
	versions, err := GetAvailableVersions()
	if err != nil {
		return "", err
	}

	featureSv, err := semver.Make("3.7.0")
	if err != nil {
		return "", err
	}
	stable, versions := versions[len(versions)-1], versions[:len(versions)-1]
	stableSv, err := semver.Make(stable)
	if err != nil {
		return "", err
	}
	// stable is the newest version of 3.6.x
	for stableSv.Compare(featureSv) >= 0 {
		stable, versions = versions[len(versions)-1], versions[:len(versions)-1]
		if stableSv, err = semver.Make(stable); err != nil {
			return "", err
		}
	}
	return stable, nil
}

// InstallPythonVersion TODO
func InstallPythonVersion(versionStr string, force bool) error {
	if ok, err := isVersionInstalled(versionStr); ok {
		if force {
			err = UninstallPythonVersion(versionStr)
			if err != nil {
				return err
			}
		} else {
			return errAlreadyInstalled
		}
	} else if err != nil {
		return err
	}

	cfg := getConfig()

	// make sure temp directory exists
	cacheDir := filepath.Join(cfg.PPrefix, versionsPath, "temp")
	if err := os.MkdirAll(cacheDir, 0700); err != nil && err != os.ErrExist {
		return err
	}

	// download the installation to that directory
	installer, err := getPythonInstaller(cfg.PMirror, versionStr, cacheDir)
	if err != nil {
		return err
	}
	logger.Debugf("installer saved at %s", installer)

	// install
	panic("not implemented")

	return nil
}

// UninstallPythonVersion TODO
func UninstallPythonVersion(versionStr string) error {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return errNotInstalled
	} else if err != nil {
		return err
	}

	panic("not implemented")
	return nil
}

// ActivatePythonVersion TODO
func ActivatePythonVersion(versionStr string) error {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return errNotInstalled
	} else if err != nil {
		return err
	}

	panic("not implemented")
	return nil
}

// Deactivate removes
func Deactivate() error {
	panic("not implemented")
	return nil
}

// InstallInfo has installation information for a given version
type InstallInfo struct {
	Executable string
	LibDir     string
	ScriptsDir string
}

// VersionFiles TODO
func VersionFiles(versionStr string) (*InstallInfo, error) {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return nil, errNotInstalled
	} else if err != nil {
		return nil, err
	}

	panic("not implemented")
	return nil, nil
}

// CallWithVersion TODO
func CallWithVersion(versionStr string, args []string) error {
	files, err := VersionFiles(versionStr)
	if err != nil {
		return err
	}
	fmt.Println(files.Executable)
	panic("not implemented")
	return nil
}

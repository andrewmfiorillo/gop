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
	defaultLoggerLevel  = loggo.WARNING
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
	// path after prefix where the active version will be stored
	activePath = "p/versions"
	// minimum allowed version
	minLegalVersion = "2.7.0"
)

// Config provides a structure for user configurable parameters
type Config struct {
	// PPrefix can be overriden by setting the P_PREFIX environment variable
	// The default is the $HOME environment variable (%HOME% on Windows).
	PPrefix string
	// PMirror can be overriden by setting the P_MIRROR environment variable
	// The default is "https://www.python.org/ftp/python/"
	PMirror string
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

// InstallInfo provides a structure for specifying the directories and executable for a given installation.
// The file paths may or may not exist, but they are always absolute.
type InstallInfo struct {
	Executable string
	BinDir     string
	LibDir     string
	IncludeDir string
	ShareDir   string
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

func getVersionDirectories(versionStr string) InstallInfo {
	cfg := getConfig()
	dirs := InstallInfo{
		Executable: "",
		BinDir:     "",
		LibDir:     "",
		ShareDir:   "",
		IncludeDir: "",
	}

	versionDir := filepath.Join(cfg.PPrefix, versionsPath, versionStr)
	if _, err := os.Stat(filepath.Join(versionDir, "bin")); err == nil {
		dirs.BinDir = filepath.Join(versionDir, "bin")
		dirs.Executable = filepath.Join(dirs.BinDir, "python")
	}

	if _, err := os.Stat(filepath.Join(versionDir, "lib")); err == nil {
		dirs.LibDir = filepath.Join(versionDir, "lib")
	}

	if _, err := os.Stat(filepath.Join(versionDir, "share")); err == nil {
		dirs.ShareDir = filepath.Join(versionDir, "share")
	}

	if _, err := os.Stat(filepath.Join(versionDir, "include")); err == nil {
		dirs.IncludeDir = filepath.Join(versionDir, "include")
	}

	return dirs
}

func getActiveDirectories() (InstallInfo, InstallInfo) {
	cfg := getConfig()
	activeDir := filepath.Join(cfg.PPrefix, activePath)

	targetDirs := InstallInfo{
		Executable: filepath.Join(activeDir, "bin", "python"),
		BinDir:     filepath.Join(activeDir, "bin"),
		LibDir:     filepath.Join(activeDir, "lib"),
		ShareDir:   filepath.Join(activeDir, "share"),
		IncludeDir: filepath.Join(activeDir, "include"),
	}
	existingDirs := InstallInfo{Executable: "", BinDir: "", LibDir: "", ShareDir: "", IncludeDir: ""}

	if _, err := os.Stat(targetDirs.BinDir); err == nil {
		existingDirs.BinDir = targetDirs.BinDir
	}

	if _, err := os.Stat(targetDirs.Executable); err == nil {
		existingDirs.Executable = targetDirs.Executable
	}

	if _, err := os.Stat(targetDirs.LibDir); err == nil {
		existingDirs.LibDir = targetDirs.LibDir
	}

	if _, err := os.Stat(targetDirs.ShareDir); err == nil {
		existingDirs.ShareDir = targetDirs.ShareDir
	}

	if _, err := os.Stat(targetDirs.IncludeDir); err == nil {
		existingDirs.IncludeDir = targetDirs.IncludeDir
	}

	return existingDirs, targetDirs
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

	installedVersions := []string{}
	for _, vDir := range versions {
		// make sure it's a directory
		if !vDir.IsDir() {
			continue
		}
		// make sure it has python in it
		// TODO: this will not work on windows (expecting "python.exe")
		if _, err := os.Stat(filepath.Join(versionsDir, vDir.Name(), "bin", "python")); os.IsNotExist(err) {
			continue
		}

		// logger.Debugf("raw version dir: %s", vDir.Name())
		vStr, err := cleanVersionString(vDir.Name())
		if err == nil {
			// logger.Debugf("clean version dir: %s", vStr)
			installedVersions = append(installedVersions, vStr)
		}
	}
	return installedVersions, nil
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

// InstallPythonVersion downloads, builds, and activates the desired version of python
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

	// make version's directory and install
	versionDir := filepath.Join(cfg.PPrefix, versionsPath, versionStr)
	if err := os.MkdirAll(cacheDir, 0700); err != nil && err != os.ErrExist {
		return err
	}
	if _, err := installPythonInstaller(installer, versionDir); err != nil {
		logger.Infof("error installing %s, deleting directory...", installer)
		_ = os.RemoveAll(versionDir)
		return err
	}

	// and remove installer
	if err := os.Remove(installer); err != nil {
		return err
	}

	return nil
}

// UninstallPythonVersion uninstalls the specified version of python
func UninstallPythonVersion(versionStr string) error {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return errNotInstalled
	} else if err != nil {
		return err
	}

	current, err := GetCurrentVersion()
	if err != nil {
		return err
	}
	if current == versionStr {
		logger.Warningf("version %s is active, deactivating...", versionStr)
		if err = Deactivate(); err != nil {
			return err
		}
	}

	cfg := getConfig()
	versionDir := filepath.Join(cfg.PPrefix, versionsPath, versionStr)
	logger.Infof("deleting %s", versionDir)
	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	return nil
}

// ActivatePythonVersion creates links to the specified version in the active directories
func ActivatePythonVersion(versionStr string) error {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return errNotInstalled
	} else if err != nil {
		return err
	}

	// remove any active directories
	if err := Deactivate(); err != nil {
		return err
	}

	// make the links
	_, activeTarget := getActiveDirectories()
	versionDirs := getVersionDirectories(versionStr)
	sources := []string{versionDirs.BinDir, versionDirs.LibDir, versionDirs.IncludeDir, versionDirs.ShareDir}
	targets := []string{activeTarget.BinDir, activeTarget.LibDir, activeTarget.IncludeDir, activeTarget.ShareDir}
	for idx := range sources {
		if err := os.Symlink(sources[idx], targets[idx]); err != nil {
			return err
		}
		logger.Infof("created link %s --> %s", sources[idx], targets[idx])
	}

	return nil
}

// Deactivate removes links for the currently active version
func Deactivate() error {
	activeExisting, _ := getActiveDirectories()
	for _, dir := range []string{activeExisting.BinDir, activeExisting.LibDir, activeExisting.IncludeDir, activeExisting.ShareDir} {
		if dir != "" {
			// it exists, remove it
			if err := os.RemoveAll(dir); err != nil {
				return err
			}
			logger.Infof("removed %s", dir)
		}
	}
	return nil
}

// VersionFiles returns the installation information for the given version
func VersionFiles(versionStr string) (*InstallInfo, error) {
	if ok, err := isVersionInstalled(versionStr); !ok {
		return nil, errNotInstalled
	} else if err != nil {
		return nil, err
	}

	files := getVersionDirectories(versionStr)
	return &files, nil
}

// CallWithVersion executes the args with the specified python version
func CallWithVersion(versionStr string, args []string) error {
	files, err := VersionFiles(versionStr)
	if err != nil {
		return err
	}
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(files.Executable, args...)
	logger.Infof("cmd: %s", files.Executable)
	logger.Infof("args: %s", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s", out)
	return nil
}

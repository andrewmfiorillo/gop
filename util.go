package pgo

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/blang/semver"
	"github.com/mholt/archiver"
)

func getPythonVersions(mirrorURL string) ([]semver.Version, error) {
	// call the mirror
	resp, err := http.Get(mirrorURL)
	if err != nil {
		return nil, err
	}

	// read the request body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	htmlStr := string(body)

	// parse the html
	versions := reIdentifier.FindAllString(htmlStr, -1)

	semVers := newSemverOrderedSet()
	for _, vstr := range versions {
		_ = semVers.Add(vstr)
	}
	return semVers.AsSlice(), nil
}

func getPythonInstaller(mirrorURL string, versionStr string, targetDir string) (string, error) {
	// TODO: check the hashes
	installerURL := getPythonInstallerURLUnix(mirrorURL, versionStr)
	if runtime.GOOS == "windows" {
		installerURL = getPythonInstallerURLWin(mirrorURL, versionStr)
	}
	filename := path.Base(installerURL)
	targetFile := filepath.Join(targetDir, filename)

	// if the file exists, we're done
	if _, err := os.Stat(targetFile); err == nil {
		logger.Infof("file exists at %s, using it...", targetFile)
		return targetFile, nil
	}

	// Create the file
	out, err := os.Create(targetFile)
	if err != nil {
		return "", err
	}
	defer out.Close()
	logger.Infof("writing to %s", targetFile)

	// Get the data
	resp, err := http.Get(installerURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	logger.Infof("got file from %s", installerURL)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return targetFile, nil
}

func installPythonInstaller(installerFile string, versionDir string) (string, error) {
	if runtime.GOOS == "windows" {
		return installPythonInstallerWin(installerFile, versionDir)
	}
	return installPythonInstallerUnix(installerFile, versionDir)
}

func getPythonInstallerURLUnix(mirrorURL string, versionStr string) string {
	return fmt.Sprintf("%s%s/Python-%s.tgz", mirrorURL, versionStr, versionStr)
}

func installPythonInstallerUnix(installerFile string, versionDir string) (string, error) {
	// the installer file is a tgz archive, so we must extract and cleanup
	if err := archiver.Unarchive(installerFile, versionDir); err != nil {
		return "", err
	}
	installerExtension := filepath.Ext(installerFile)
	installerFilename := filepath.Base(installerFile)
	installerFilestem := installerFilename[0 : len(installerFilename)-len(installerExtension)]
	extractedDir := filepath.Join(versionDir, installerFilestem)
	srcDir := filepath.Join(versionDir, "src")
	if err := os.Rename(extractedDir, srcDir); err != nil {
		return "", err
	}
	logger.Debugf("extracted to %s", srcDir)

	// now we configure and build

	// ./configure --prefix="$dir"
	logger.Infof("running `./configure --prefix=%s`", versionDir)
	cmd := exec.Command("./configure", fmt.Sprintf("--prefix=%s", versionDir))
	cmd.Dir = srcDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Debugf("./configure output: %s", out)
		logger.Debugf("`./configure` error: %s", err)
		return "", fmt.Errorf("unable to configure python source in %s", srcDir)
	}

	// make &> /dev/null
	logger.Infof("running `make`")
	cmd = exec.Command("make")
	cmd.Dir = srcDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		logger.Debugf("`./make output`: %s", out)
		logger.Debugf("`make` error: %s", err)
		return "", fmt.Errorf("unable to make python source in %s", srcDir)
	}

	// make install &> /dev/null
	logger.Infof("running `make install`")
	cmd = exec.Command("make", "install")
	cmd.Dir = srcDir
	out, err = cmd.CombinedOutput()
	if err != nil {
		logger.Debugf("`make install` output: %s", out)
		logger.Debugf("`make install` error: %s", err)
		return "", fmt.Errorf("unable to make install python source in %s", srcDir)
	}

	// make links
	python3Path := filepath.Join(versionDir, "bin", "python3")
	pythonPath := filepath.Join(versionDir, "bin", "python")
	if _, err = os.Stat(python3Path); err == nil {
		if err = os.Symlink(python3Path, pythonPath); err != nil {
			return "", err
		}
	}

	pip3Path := filepath.Join(versionDir, "bin", "pip3")
	pipPath := filepath.Join(versionDir, "bin", "pip")
	if _, err = os.Stat(pipPath); os.IsNotExist(err) {
		if err = os.Symlink(pip3Path, pipPath); err != nil {
			return "", err
		}
	}

	// cleanup src directory
	if err := os.RemoveAll(srcDir); err != nil {
		return "", err
	}

	// try checking its version
	vStr, err := getPythonBinVersion(pythonPath)
	if err != nil {
		return "", err
	} else if vStr != filepath.Base(versionDir) {
		return "", fmt.Errorf("installed python version %s mismatches specified", vStr)
	}

	return pythonPath, nil
}

func getPythonInstallerURLWin(mirrorURL string, versionStr string) string {
	return fmt.Sprintf("%s%s/python-%s.amd64.msi", mirrorURL, versionStr, versionStr)
}

func installPythonInstallerWin(installerFile string, versionDir string) (string, error) {
	panic("not implemented")
	return "", nil
}

func getPythonBinVersion(pythonExec string) (string, error) {
	out, err := exec.Command(pythonExec, "--version").Output()
	if err != nil {
		return "", err
	}
	return cleanVersionString(string(out))
}

func checkConfiguration(cfg Config) error {
	_, dirs := getActiveDirectories()

	// make sure the bin directory is on PATH
	dirsOnPath := filepath.SplitList(os.Getenv("PATH"))
	binDir := dirs.BinDir
	isOnPath := false
	for _, onPath := range dirsOnPath {
		if onPath == binDir {
			isOnPath = true
		}
	}
	if !isOnPath {
		logger.Warningf("bin directory `%s` is not on PATH", binDir)
	}

	return nil
}

func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type semverOrderedSet struct {
	versions []semver.Version
	counts   map[string]bool
}

func newSemverOrderedSet() *semverOrderedSet {
	return &semverOrderedSet{
		versions: make([]semver.Version, 0),
		counts:   make(map[string]bool),
	}
}

func (sv *semverOrderedSet) Len() int {
	return len(sv.versions)
}

func (sv *semverOrderedSet) Less(i, j int) bool {
	return sv.versions[i].Compare(sv.versions[j]) <= -1
}

func (sv *semverOrderedSet) Swap(i, j int) {
	sv.versions[i], sv.versions[j] = sv.versions[j], sv.versions[i]
}

func (sv *semverOrderedSet) Add(versionStr string) error {
	sver, err := semver.Make(versionStr)
	if err != nil {
		return err
	}
	if _, ok := sv.counts[sver.String()]; !ok {
		sv.counts[sver.String()] = true
		sv.versions = append(sv.versions, sver)
	}
	return nil
}

func (sv *semverOrderedSet) AsSlice() []semver.Version {
	sort.Sort(sv)
	slice := make([]semver.Version, 0, len(sv.versions))
	for _, sver := range sv.versions {
		slice = append(slice, sver)
	}
	return slice
}

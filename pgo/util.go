package pgo

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/blang/semver"
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

	installerURL := getPythonInstallerURLUnix(mirrorURL, versionStr)
	if runtime.GOOS == "windows" {
		installerURL = getPythonInstallerURLWin(mirrorURL, versionStr)
	}
	filename := path.Base(installerURL)
	targetFile := filepath.Join(targetDir, filename)

	// Create the file
	out, err := os.Create(targetFile)
	if err != nil {
		return "", err
	}
	defer out.Close()
	logger.Debugf("writing to %s", targetFile)

	// Get the data
	logger.Debugf("getting file at %s", installerURL)
	resp, err := http.Get(installerURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return targetFile, nil
}

func getPythonInstallerURLUnix(mirrorURL string, versionStr string) string {
	return fmt.Sprintf("%s%s/Python-%s.tgz", mirrorURL, versionStr, versionStr)
}

func getPythonInstallerURLWin(mirrorURL string, versionStr string) string {
	return fmt.Sprintf("%s%s/python-%s.amd64.msi", mirrorURL, versionStr, versionStr)
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

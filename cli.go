package pgo

import (
	"fmt"

	"github.com/juju/loggo"
	"github.com/urfave/cli"
)

// MakeApp constructs a configured CLI application
func MakeApp() *cli.App {
	logger.SetLogLevel(defaultLoggerLevel)
	app := cli.NewApp()
	app.Name = "p"
	app.Usage = "get Going with Python version management"
	app.Version = "0.0.1"
	app.Before = func(c *cli.Context) error {
		cfg := getConfig()
		if err := checkConfiguration(cfg); err != nil {
			return err
		}

		if c.Bool("verbose") {
			logger.SetLogLevel(loggo.INFO)
		}
		return nil
	}
	app.Action = ActivateVersion
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "verbose"},
	}
	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"list"},
			Usage:   "Output the versions of Python available",
			Action:  ListAvailable,
			Subcommands: []cli.Command{
				{
					Name:     "installed",
					HelpName: "ls installed",
					Usage:    "Output the installed versions of Python",
					Action:   ListInstalled,
				},
				{
					Name:     "latest",
					HelpName: "ls latest",
					Usage:    "Output the latest Python version available",
					Action:   ShowLatest,
				},
				{
					Name:     "stable",
					HelpName: "ls stable",
					Usage:    "Output the latest stable Python version available",
					Action:   ShowStable,
				},
			},
		},
		{
			Name:   "latest",
			Usage:  "Activate to the latest Python release",
			Action: ActivateLatest,
		},
		{
			Name:   "stable",
			Usage:  "Activate to the latest stable Python release",
			Action: ActivateStable,
		},
		{
			Name:   "status",
			Usage:  "Output current status",
			Action: ShowStatus,
		},
		{
			Name:      "install",
			Usage:     "Install Python <version> but do NOT activate",
			ArgsUsage: "<version> --force",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "force"},
			},
			Action: InstallVersion,
		},
		{
			Name:      "use",
			Usage:     "Execute Python <version> with [args ...]",
			ArgsUsage: "<version> [args ...]",
			Action:    UseVersion,
		},
		{
			Name:      "bin",
			Usage:     "Output bin path for <version>",
			ArgsUsage: "<version>",
			Action:    ShowVersion,
		},
		{
			Name:      "rm",
			Usage:     "Remove the given version(s)",
			ArgsUsage: "<version ...>",
			Action:    RemoveVersion,
		},
		{
			Name:    "default",
			Aliases: []string{"disable"},
			Usage:   "Use default (system) Python installation",
			Action:  ActivateDefault,
		},
	}
	cli.AppHelpTemplate = `
Name:
{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

Commands:
    p <version>{{ "\t" }}Activate to Python <version>{{range .Commands}}
    p {{join .Names ", "}} {{ .ArgsUsage }}{{ "\t"}}{{.Usage}}{{ if .Subcommands }}{{range .Subcommands}}{{ "\n        " }}p {{ .HelpName }} {{ .ArgsUsage }}{{ "\t"}}{{.Usage}}{{end}}{{end}}{{end}}

Options:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}

`
	return app
}

var errNoVersionString = fmt.Errorf("no version string given")

func getVersionString(c *cli.Context) (string, error) {
	if !c.Args().Present() {
		return "", errNoVersionString
	}
	vstr := c.Args().First()

	return cleanVersionString(vstr)
}

// ListAvailable .
func ListAvailable(c *cli.Context) error {
	versions, err := GetAvailableVersions()
	if err != nil {
		return err
	}
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return err
	}

	for _, vStr := range versions {
		if vStr == currentVersion {
			fmt.Printf("--> %s \n", vStr)
		} else {
			fmt.Printf("    %s \n", vStr)
		}
	}
	return nil
}

// ListInstalled .
func ListInstalled(c *cli.Context) error {
	versions, err := GetInstalledVersions()
	if err != nil {
		return err
	}
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		return err
	}

	for _, vStr := range versions {
		if vStr == currentVersion {
			fmt.Printf("--> %s \n", vStr)
		} else {
			fmt.Printf("    %s \n", vStr)
		}
	}
	return nil
}

// ShowLatest .
func ShowLatest(c *cli.Context) error {
	latest, err := GetLatestVersion()
	if err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

// ShowStable .
func ShowStable(c *cli.Context) error {
	stable, err := GetStableVersion()
	if err != nil {
		return err
	}
	fmt.Println(stable)
	return nil
}

// ShowStatus .
func ShowStatus(c *cli.Context) error {
	vstr, err := GetCurrentVersion()
	if err != nil {
		return err
	}
	fmt.Println("current version:", vstr)
	return nil
}

// ActivateLatest installs (if necessary) and activates the latest available version of python
func ActivateLatest(c *cli.Context) error {
	latest, err := GetLatestVersion()
	if err != nil {
		return err
	}

	isInstalled, err := isVersionInstalled(latest)
	if err != nil {
		return err
	}
	if !isInstalled {
		if err := InstallPythonVersion(latest, false); err != nil {
			return err
		}
	}

	if err = ActivatePythonVersion(latest); err != nil {
		return err
	}
	fmt.Println(latest)
	return nil
}

// ActivateStable installs (if necessary) and activates the latest stable version of python
func ActivateStable(c *cli.Context) error {
	stable, err := GetStableVersion()
	if err != nil {
		return err
	}

	isInstalled, err := isVersionInstalled(stable)
	if err != nil {
		return err
	}
	if !isInstalled {
		if err := InstallPythonVersion(stable, false); err != nil {
			return err
		}
	}

	if err = ActivatePythonVersion(stable); err != nil {
		return err
	}
	fmt.Println(stable)
	return nil
}

// ActivateVersion installs and activates the given version of python
func ActivateVersion(c *cli.Context) error {
	// get version string
	vstr, err := getVersionString(c)
	if err == errNoVersionString {
		return ListInstalled(c)
	} else if err != nil {
		return err
	}
	logger.Debugf("specified version: %s", vstr)

	// is is installed?
	installedVersions, err := GetInstalledVersions()
	if err != nil {
		return err
	}
	if !stringContains(installedVersions, vstr) {
		logger.Infof("version %s not installed, installing...", vstr)
		if err = InstallPythonVersion(vstr, false); err != nil {
			return err
		}
	}

	// activate it
	if err := ActivatePythonVersion(vstr); err != nil {
		return err
	}
	fmt.Println("activated", vstr)
	return nil
}

// InstallVersion installs the specified version of python but does not activate
func InstallVersion(c *cli.Context) error {
	// get version string
	vstr, err := getVersionString(c)
	if err != nil {
		return err
	}
	logger.Debugf("specified version: %s", vstr)

	force := c.Bool("force")
	if err = InstallPythonVersion(vstr, force); err != nil {
		return err
	}
	fmt.Println(vstr)
	return nil
}

// UseVersion executes a command with the given arguments
func UseVersion(c *cli.Context) error {
	// get version string
	vstr, err := getVersionString(c)
	if err != nil {
		return err
	}
	logger.Debugf("specified version: %s", vstr)
	return CallWithVersion(vstr, c.Args().Tail())
}

// ShowVersion displays the path to the specified version of python
func ShowVersion(c *cli.Context) error {
	// get version string
	vstr, err := getVersionString(c)
	if err != nil {
		return err
	}
	logger.Debugf("specified version: %s", vstr)

	files, err := VersionFiles(vstr)
	if err != nil {
		return err
	}
	fmt.Println(files.Executable)
	return nil
}

// RemoveVersion uninstalls the specified version
func RemoveVersion(c *cli.Context) error {
	// get version string
	vstr, err := getVersionString(c)
	if err != nil {
		return err
	}
	logger.Debugf("specified version: %s", vstr)
	if err := UninstallPythonVersion(vstr); err != nil {
		return err
	}
	fmt.Println("uninstalled %s", vstr)
	return nil
}

// ActivateDefault reverts the to default sytem python
func ActivateDefault(c *cli.Context) error {
	if err := Deactivate(); err != nil {
		return err
	}
	vstr, err := GetCurrentVersion()
	if err != nil {
		logger.Errorf("no system python installed!")
		return err
	}
	fmt.Println("system python: %s", vstr)
	return nil
}

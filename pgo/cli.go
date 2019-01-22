package pgo

import (
	"fmt"

	"github.com/urfave/cli"
)

// MakeApp constructs a configured CLI application
func MakeApp() *cli.App {
	app := cli.NewApp()
	app.Name = "p"
	app.Usage = "get Going with Python version management"
	app.Version = "0.0.1"
	app.Action = ActivateVersion
	app.Commands = []cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"list"},
			Usage:   "Output the versions of Python available",
			Action:  ListAvailable,
			Subcommands: []cli.Command{
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
			ArgsUsage: "<version>",
			Action:    InstallVersion,
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
	return nil
}

// ListInstalled .
func ListInstalled(c *cli.Context) error {
	return nil
}

// ShowLatest .
func ShowLatest(c *cli.Context) error {
	return nil
}

// ShowStable .
func ShowStable(c *cli.Context) error {
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

// ActivateLatest .
func ActivateLatest(c *cli.Context) error {
	return nil
}

// ActivateStable .
func ActivateStable(c *cli.Context) error {
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

	// activate it

	return nil
}

// InstallVersion .
func InstallVersion(c *cli.Context) error {
	return nil
}

// UseVersion .
func UseVersion(c *cli.Context) error {
	return nil
}

// ShowVersion .
func ShowVersion(c *cli.Context) error {
	return nil
}

// RemoveVersion .
func RemoveVersion(c *cli.Context) error {
	return nil
}

// ActivateDefault .
func ActivateDefault(c *cli.Context) error {
	return nil
}

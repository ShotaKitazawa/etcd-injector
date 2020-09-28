package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	appName    = "etcd-injector"
	appVersion = "dummy"
)

type config struct {
	AppName       string
	AppVersion    string
	SrcEndpoints  []string
	SrcUsername   string
	SrcPassword   string
	SrcDirectory  string
	DstEndpoints  []string
	DstUsername   string
	DstPassword   string
	DstDirectory  string
	RulesFilepath string
	IgnoreKeys    []string
	DeleteEnable  bool
	LoggingEnable bool
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Usage = "recursive copy & inject json value of etcd"
	app.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "src-endpoints",
			Usage:    "source endpoints of etcd",
			Required: true,
			EnvVars:  []string{"ETCD_SRC_ENDPOINTS"},
		},
		&cli.StringFlag{
			Name:    "src-username",
			Usage:   "username of source etcd",
			EnvVars: []string{"ETCD_SRC_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "src-password",
			Usage:   "password of source etcd",
			EnvVars: []string{"ETCD_SRC_PASSWORD"},
		},
		&cli.StringFlag{
			Name:     "src-directory",
			Aliases:  []string{"s"},
			Usage:    "source directory of etcd",
			Required: true,
			EnvVars:  []string{"ETCD_SRC_DIRECTORY"},
		},
		// TODO
		// &cli.StringFlag{
		//	Name:     "src-cafile",
		// },
		// &cli.StringFlag{
		//	Name:     "src-certfile",
		// },
		// &cli.StringFlag{
		//	Name:     "src-keyfile",
		// },
		&cli.StringSliceFlag{
			Name:     "dst-endpoints",
			Usage:    "destination endpoints of etcd",
			Required: true,
			EnvVars:  []string{"ETCD_DST_ENDPOINTS"},
		},
		&cli.StringFlag{
			Name:    "dst-username",
			Usage:   "username of destination etcd",
			EnvVars: []string{"ETCD_DST_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "dst-password",
			Usage:   "password of destination etcd",
			EnvVars: []string{"ETCD_DST_PASSWORD"},
		},
		&cli.StringFlag{
			Name:     "dst-directory",
			Aliases:  []string{"d"},
			Usage:    "destination directory of etcd",
			Required: true,
			EnvVars:  []string{"ETCD_DST_DIRECTORY"},
		},
		// TODO
		// &cli.StringFlag{
		//	Name:     "dst-cafile",
		// },
		// &cli.StringFlag{
		//	Name:     "dst-certfile",
		// },
		// &cli.StringFlag{
		//	Name:     "dst-keyfile",
		// },
		&cli.StringFlag{
			Name:     "rules-filepath",
			Aliases:  []string{"f"},
			Usage:    "path of file written injection rules",
			Required: true,
			EnvVars:  []string{"RULES_FILEPATH"},
		},
		&cli.StringSliceFlag{
			Name:    "ignore",
			Usage:   `specified "--ignore=/key", "xxx" is excluded from copy target`,
			EnvVars: []string{"IGNORE_KEYS"},
		},
		&cli.BoolFlag{
			Name:  "delete",
			Usage: `delete dst key if does not exist in src (like "rsync --delete")`,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"x"},
			Usage:   "output results of replacement",
		},
	}
	app.Action = func(c *cli.Context) error {
		return Run(config{
			AppName:       appName,
			AppVersion:    appVersion,
			SrcEndpoints:  c.StringSlice("src-endpoints"),
			SrcUsername:   c.String("src-username"),
			SrcPassword:   c.String("src-password"),
			SrcDirectory:  c.String("src-directory"),
			DstEndpoints:  c.StringSlice("dst-endpoints"),
			DstUsername:   c.String("dst-username"),
			DstPassword:   c.String("dst-password"),
			DstDirectory:  c.String("dst-directory"),
			RulesFilepath: c.String("rules-filepath"),
			IgnoreKeys:    c.StringSlice("ignore"),
			DeleteEnable:  c.Bool("delete"),
			LoggingEnable: c.Bool("verbose"),
		})
	}
	app.EnableBashCompletion = true
	return app
}

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

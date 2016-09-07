package app

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"shadowsocks-go/pkg/permissions"

	"github.com/golang/glog"
	"github.com/urfave/cli"
)

func Run() {
	app := cli.NewApp()
	registerCmd(app)
	app.Run(os.Args)
}

func registerCmd(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "GenerateLicense",
			Aliases: []string{"gen"},
			Usage:   "generate license by hw code",
			Action: func(c *cli.Context) error {
				// hwcode, err := permissions.GenerateHardWareCode()
				hwcode, err := hex.DecodeString(c.String("hardwareCode"))
				license, err := permissions.GenerateLicense(hwcode)
				ioutil.WriteFile(".license", []byte(license), os.FileMode(0400))
				glog.Infof("Get license %v \r\v", license)
				return err
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "hardwareCode, hwcode",
					Usage: "hardware code for generate license",
				},
			},
		},
		{
			Name:    "GenerateHWCode",
			Aliases: []string{"gencode"},
			Usage:   "gen code",
			Action: func(c *cli.Context) error {
				hwcode, err := permissions.GenerateHardWareCode()
				glog.Infof("get code %v err %v", hwcode, err)
				return nil
			},
		},
		{
			Name:    "CheckLicense",
			Aliases: []string{"check"},
			Usage:   "check license",
			Action: func(c *cli.Context) error {
				result := permissions.PermissionsCheck(c.String("license"))
				glog.Infof("check license result %v", result)
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "license",
					Usage: "check input license",
				},
			},
		},
	}
}

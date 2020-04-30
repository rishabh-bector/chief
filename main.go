package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	access "bector.dev/chief/access"
	config "bector.dev/chief/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	app := &cli.App{
		Name:  "chief",
		Usage: "a continuous integration/deployment server for hobbyists",
		Commands: []*cli.Command{
			{
				Name:   "setup",
				Usage:  "set up the chief server",
				Action: setupServer,
			},
			{
				Name:  "access",
				Usage: "manage user access",
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "'chief access add <username>' to add a new user",
						Action: access.AddUserCommand,
					},
					{
						Name:   "remove",
						Usage:  "remove an existing user",
						Action: access.RemoveUser,
					},
					{
						Name:   "list",
						Usage:  "list all users & clearances",
						Action: access.List,
					},
				},
			},
			{
				Name:  "pipeline",
				Usage: "manage pipelines",
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "add a new pipeline to the chief server",
						Action: access.AddUserCommand,
					},
				},
			},
			{
				Name:   "start",
				Usage:  "start the chief server",
				Action: startServer,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "detach",
						Value: false,
						Usage: "start the chief server and detach, running it in the background",
					},
				},
			},
			{
				Name:   "status",
				Usage:  "get the chief server status",
				Action: serverStatus,
			},
			{
				Name:   "kill",
				Usage:  "end the chief server process",
				Action: killServer,
			},
			{
				Name:   "uninstall",
				Usage:  "uninstall the chief server",
				Action: access.Uninstall,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupServer(c *cli.Context) error {
	log.Debug("Obtaining user's home dir...")
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	log.Debugf("Searching for .chief dir in %s...", home)
	files, err := ioutil.ReadDir(home)
	if err != nil {
		return err
	}

	found := false
	for _, f := range files {
		if f.Name() == ".chief" {
			found = true
			log.Debug("Found existing .chief dir!")
		}
	}

	if !found {
		log.Debug("Existing .chief dir not found, creating...")
		err := os.Mkdir(fmt.Sprintf("%s/%s", home, ".chief"), 0755)
		if err != nil {
			return err
		} else {
			log.Debugf("Created .chief dir in %s!", home)
		}
	}

	log.Debugf("Searching for config.json in %s/.chief...", home)
	chiefFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", home, ".chief"))
	if err != nil {
		return err
	}

	foundChief := false
	for _, f := range chiefFiles {
		if f.Name() == "config.json" {
			foundChief = true
			log.Debug("Found existing config.json!")
			err = config.LoadFromDisk()
			if err != nil {
				return err
			}
		}
	}

	if !foundChief {
		log.Debug("Existing config.json not found, creating...")
		err = config.Setup()
		if err != nil {
			return err
		}
	}

	// Access setup
	err = access.Setup()
	if err != nil {
		return err
	}

	return nil
}

func startServer(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}

	s := NewChiefServer(config.Global())

	if c.Bool("detach") {
		cmd := exec.Command("chief", "start", "&")
		err := cmd.Start()
		return err
	}

	err := s.Start()
	return err
}

func serverStatus(c *cli.Context) error {
	resp, err := http.Get("http://localhost:2222/status")
	if err != nil {
		fmt.Printf(StatusMessage, STATUS_STOPPED)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var status ServerStatus
	json.Unmarshal(body, &status)

	fmt.Printf(StatusMessage, status.Status)
	return nil
}

func killServer(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}

	if err := access.Ensure(config.MASTER_CLEARANCE); err != nil {
		return err
	}

	http.Get("http://localhost:2222/kill")
	return nil
}

const StatusMessage = `
--------------------------------------------------
| Chief Status: %s
--------------------------------------------------

`

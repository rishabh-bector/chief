package access

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"

	config "bector.dev/chief/config"
)

var Clearances = map[string]int{
	config.MASTER_CLEARANCE: 2,
	config.NORMAL_CLEARANCE: 1,
}

// Setup is run on installation to add the master user
func Setup() error {
	if err := config.Ensure(); err != nil {
		return err
	}

	// Check for existing master user
	for _, u := range config.Global().Access {
		if u.Clearance == config.MASTER_CLEARANCE {
			log.Debug("Found existing master user!")
			return nil
		}
	}

	getInput(setupMessage)

	usr := getInput("Enter username: ")
	err := AddUser(usr, config.MASTER_CLEARANCE)
	if err != nil {
		return err
	}
	return nil
}

// AddUser adds a new user
func AddUser(user string, clearance string) error {
	if user == "" {
		return errors.New("username cannot be blank")
	}

	pass := getInput("Enter new user's password: ")
	if pass == "" {
		return errors.New("password cannot be blank")
	}

	rUser := strings.ReplaceAll(user, "\n", "")
	rPass := strings.ReplaceAll(pass, "\n", "")

	hsh, err := HashPassword(rPass)
	if err != nil {
		return err
	}

	usr := config.User{
		PassHash:  hsh,
		Clearance: clearance,
	}

	if _, ok := config.Global().Access[rUser]; ok {
		return errors.New("user already exists")
	}

	config.Global().Access[rUser] = usr
	config.WriteToDisk()

	return nil
}

// Ensure prompts the user to login, and verifies their clearance
func Ensure(clearance string) error {
	fmt.Printf(loginMessage, config.MASTER_CLEARANCE)

	var usr string

	if len(config.Global().Access) > 1 {
		usr = getInput("Enter username: ")
	} else {
		for k, _ := range config.Global().Access {
			usr = k
		}
	}

	if u, ok := config.Global().Access[usr]; ok {
		if Clearances[u.Clearance] < Clearances[clearance] {
			return fmt.Errorf("failed, user is of %s clearance, minimum %s clearance required", u.Clearance, clearance)
		}

		// Get password
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		if CheckPasswordHash(string(bytePassword), u.PassHash) {
			fmt.Println()
			return nil
		} else {
			return errors.New("invalid password")
		}
	}
	return fmt.Errorf("user '%s' does not exist", usr)
}

// AddUserCommand is called by CLI
func AddUserCommand(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}
	if err := Ensure(config.MASTER_CLEARANCE); err != nil {
		return err
	}

	return AddUser(c.Args().Get(0), config.NORMAL_CLEARANCE)
}

// RemoveUser removes a user
func RemoveUser(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}

	if err := Ensure(config.MASTER_CLEARANCE); err != nil {
		return err
	}

	usr := c.Args().Get(0)
	if u, ok := config.Global().Access[usr]; ok {
		if u.Clearance == config.MASTER_CLEARANCE {
			return errors.New("cannot remove a user with master clearance")
		}
		delete(config.Global().Access, usr)
		config.WriteToDisk()
		return nil
	}
	return fmt.Errorf("user '%s' does not exist", usr)
}

// List lists all users & their clearances
func List(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}

	if err := Ensure(config.MASTER_CLEARANCE); err != nil {
		return err
	}

	fmt.Print(listMessage)
	for uname, u := range config.Global().Access {
		fmt.Printf("%s: %s\n", uname, u.Clearance)
	}
	fmt.Println()

	return nil
}

// Uninstall stops all pipelines and removes chief from the system
func Uninstall(c *cli.Context) error {
	if err := config.Ensure(); err != nil {
		return err
	}

	conf := getInput(deleteMessage)
	if conf != "yes" {
		log.Warn("Aborting uninstallation!")
		return nil
	}
	log.Warn("Proceeding with uninstallation...")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.Remove(fmt.Sprintf("%s/.chief/config.json", home))
	if err != nil {
		return err
	}

	err = os.Remove(fmt.Sprintf("%s/.chief", home))
	if err != nil {
		return err
	}

	return nil
}

func getInput(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	input, _ := reader.ReadString('\n')
	return strings.ReplaceAll(input, "\n", "")
}

const loginMessage = `

--------------------------------------------------
| Chief Access Control | Clearance: %s |
--------------------------------------------------

`

const setupMessage = `

--------------------------------------------------
| Chief access setup
--------------------------------------------------

You will now be prompted to create a master username/
password for use with the Chief server. Press
<return> to continue.
`

const deleteMessage = `

--------------------------------------------------
| Uninstalling Chief 
--------------------------------------------------

WARNING: You are about to uninstall the Chief server
from your system. All pipelines will cease immediately,
and all files will be cleaned from the system.

Continue? (yes/no): `

const listMessage = `

user: clearance 
-------------------------
`

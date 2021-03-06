package me

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/runabove/sail/internal"
)

var cmdMeSetPassword = &cobra.Command{
	Use:   "setPassword",
	Short: "Set account password: sail me setPassword [<password>]",
	Long: `Set account password: sail me setAcl <ip> [<ip> ... ]
	"example: sail me setPassword"
	Note: After running this command, you will need to run "docker login sailabove.io"
	`,
	Aliases: []string{"password", "set-password"},
	Run:     cmdSetPassword,
}

// users struct holds all parameters sent to /users
type users struct {
	Password string `json:"password,omitempty"`
}

func cmdSetPassword(cmd *cobra.Command, args []string) {
	var cmdUsersBody users

	switch len(args) {
	case 1:
		cmdUsersBody.Password = args[0]
	case 0:
		fmt.Fprint(os.Stderr, "Password: ")
		password := gopass.GetPasswd()

		fmt.Fprint(os.Stderr, "Confirm password: ")
		confirm := gopass.GetPasswd()

		if !bytes.Equal(password, confirm) {
			fmt.Fprintln(os.Stderr, "Error: Passwords do not match")
			return
		}

		if len(password) == 0 {
			fmt.Fprintln(os.Stderr, "Error: Password Required")
			return
		}

		cmdUsersBody.Password = string(password[:])
	default:
		fmt.Fprintln(os.Stderr, "Invalid usage. sail me password [<password>]. Please see sail me password --help")
		return
	}

	jsonStr, err := json.Marshal(cmdUsersBody)
	internal.Check(err)
	internal.FormatOutputDef(internal.ReqWantJSON("PUT", http.StatusOK, "/users", jsonStr))
}

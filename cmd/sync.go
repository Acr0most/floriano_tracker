package cmd

import (
	"fmt"

	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"tracker/helpers"
	"tracker/tracker"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Get the frames from the server and push the new ones.",
	Long: `Get the frames from the server and push the new ones.

  The URL of the server and the User Token must be defined via the 'watson
	config' command.

  Example:

  $ tracker config backend.url http://localhost:4242
  $ tracker config backend.token 7e329263e329
  $ tracker sync
  Received 42 frames from the server
  Pushed 23 frames to the server`,
	Run: sync,
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func sync(cmd *cobra.Command, args []string) {
	url := viper.GetString("backend.url")
	token := viper.GetString("backend.token")

	if url == "" || token == "" {
		fmt.Printf("Error: %s\n", helpers.PrintRed("You need to set backend url and token before being able to sync."))
		fmt.Println("        tracker config backend.url http://some.url")
		fmt.Println("        tracker config backend.token mytoken")
		return
	}

	frames := tracker.GetFrames()
	b, err := json.Marshal(frames)

	if err != nil {
		fmt.Printf("Error: %s. %s\n", helpers.PrintRed("Unabled to marshal frames"), err.Error())
		return
	}

	err = upload(url, token, b)
	if err != nil {
		fmt.Printf("Error: %s\n", helpers.PrintRed(err.Error()))
		return
	}

	fmt.Println("Sync successul")
}

func upload(url, token string, data []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-Token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		return errors.New(resp.Status)
	}

	return nil
}

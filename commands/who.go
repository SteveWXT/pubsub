package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/SteveWXT/pubsub/clients"
)

var (
	whoCmd = &cobra.Command{
		Hidden:        true,
		Use:           "who",
		Short:         "List connection stats",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: who,
	}
)

// init
func init() {
	whoCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
}

// who gets connection stats for a mist server
func who(ccmd *cobra.Command, args []string) error {

	// create new mist client
	client, err := clients.New(host)
	if err != nil {
		fmt.Printf("Failed to connect to '%s' - %s\n", host, err.Error())
		return err
	}

	// who related
	err = client.Who()
	if err != nil {
		fmt.Printf("Failed to who - %s\n", err.Error())
		return err
	}

	msg := <-client.Messages()
	fmt.Println(msg.Data)

	return nil
}

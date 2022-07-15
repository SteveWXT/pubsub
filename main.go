package main

import (
	"fmt"

	"github.com/SteveWXT/pubsub/commands"
)

func main() {
	if err := commands.PubSubCmd.Execute(); err != nil && err.Error() != "" {
		fmt.Println(err.Error())
	}
}

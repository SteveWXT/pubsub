// Package commands ...
package commands

import (
	"fmt"
	"path/filepath"

	"github.com/SteveWXT/pubsub/server"
	"github.com/jcelliott/lumber"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	host = "127.0.0.1:1445" // host clients will connect to
	tags []string           // tags to publish and [un]subscribe to/from

	config   string // location of the config file
	showVers bool   // whether to show version info and exit or not

	// to be populated by linker
	version string
	commit  string

	// PubSubCmd ...
	PubSubCmd = &cobra.Command{
		Use:           "pubsub",
		Short:         "PubSub is a simple pub/sub service for tagged messages",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		// parse the config if one is provided, or use the defaults
		PersistentPreRunE: readConfig,

		// print version or help, or continue, depending on flag settings
		PreRunE: preFlight,

		// either run as a server, or run as a CLI depending on what flags
		// are provided
		RunE: start,
	}
)

func readConfig(ccmd *cobra.Command, args []string) error {
	// if --version is passed print the version info
	if showVers {
		fmt.Printf("PubSub %s (%s)\n", version, commit)
		return fmt.Errorf("")
	}

	// if --config is passed, attempt to parse the config file
	if config != "" {
		filename := filepath.Base(config)
		viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
		viper.AddConfigPath(filepath.Dir(config))

		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("Failed to read config file - %s", err.Error())
		}
	}

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	// if --server is not passed, print help
	if !viper.GetBool("server") {
		ccmd.HelpFunc()(ccmd, args)
		return fmt.Errorf("") // no error, just exit
	}

	return nil
}

func start(ccmd *cobra.Command, args []string) error {
	// configure the logger
	lumber.Prefix("[pubsub]")
	lumber.Level(lumber.LvlInt(viper.GetString("log-level")))

	if err := server.Start(viper.GetStringSlice("listeners")); err != nil {
		return fmt.Errorf("One or more servers failed to start - %s", err.Error())
	}

	return nil
}

func init() {

	// persistent config flags
	PubSubCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	viper.BindPFlag("log-level", PubSubCmd.PersistentFlags().Lookup("log-level"))

	PubSubCmd.Flags().StringSlice("listeners", []string{"tcp://127.0.0.1:1445", "ws://127.0.0.1:8888"}, "A comma delimited list of servers to start")
	viper.BindPFlag("listeners", PubSubCmd.Flags().Lookup("listeners")) // no reason to have "http://127.0.0.1:8080" too, it only has /ping

	PubSubCmd.Flags().StringVar(&config, "config", config, "Path to config file")
	viper.BindPFlag("config", PubSubCmd.Flags().Lookup("config"))

	PubSubCmd.Flags().Bool("server", false, "Run PubSub as a server")
	viper.BindPFlag("server", PubSubCmd.Flags().Lookup("server"))

	PubSubCmd.Flags().BoolVarP(&showVers, "version", "v", false, "Display the current version of this CLI")

	// commands
	PubSubCmd.AddCommand(pingCmd)
	PubSubCmd.AddCommand(subscribeCmd)
	PubSubCmd.AddCommand(publishCmd)

	// hidden/aliased commands
	PubSubCmd.AddCommand(listCmd)
	PubSubCmd.AddCommand(whoCmd)
	PubSubCmd.AddCommand(messageCmd)
	PubSubCmd.AddCommand(sendCmd)
}

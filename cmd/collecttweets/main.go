package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// Create root command.
	rootCmd := NewCommand()

	// Execute adds all child commands to the root command and sets flags appropriately.
	if err := rootCmd.Execute(); err != nil {
		// When converting error to string, preserve new lines.
		msg := strings.Replace(fmt.Sprintf("%v", err), `\n`, "\n", -1)
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", msg)
		os.Exit(1)
	}
}

const (
	vBearerToken = "bearer-token"
	vTopic       = "topic"

	flagBearerToken = "bearer-token"
	flagTopic       = "topic"

	shortHandTopic       = "t"
	shortHandBearerToken = "b"
)

// NewCommand returns a new cobra.Command for the base command
func NewCommand() *cobra.Command {
	viper.SetEnvPrefix("twitter")

	cmd := &cobra.Command{
		Use:   "collecttweets -b <bearer-token> -t <topic>",
		Short: "collecttweets is a tool for collecting tweets for a topic.",
		Long:  `collecttweets is a tool for collecting tweets for a topic from Twitter.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.CommandLine.Parse([]string{}); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract the arguments.
			bearerToken := viper.GetString(vBearerToken)
			topic := viper.GetString(vTopic)
			if bearerToken == "" || topic == "" {
				return cmd.Usage()
			}

			return collectTweets(bearerToken, topic)
		},
	}

	cmd.PersistentFlags().CountP("verbose", "v", "Increase verbosity. May be given multiple times.")

	cmd.Flags().StringP(flagBearerToken, shortHandBearerToken, "", "Bearer token")
	_ = viper.BindPFlag(vBearerToken, cmd.Flags().Lookup(flagBearerToken))

	cmd.Flags().StringP(flagTopic, shortHandTopic, "", "Topic, e.g. Facebook")
	_ = viper.BindPFlag(vTopic, cmd.Flags().Lookup(flagTopic))

	return cmd
}

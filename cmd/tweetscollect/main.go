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
	vOutputFile  = "output-file"
	vTopic       = "topic"

	flagBearerToken = "bearer-token"
	flagOutputFile  = "output-file"
	flagTopic       = "topic"

	shortHandBearerToken = "b"
	shortHandOutputFile  = "o"
	shortHandTopic       = "t"
)

// NewCommand returns a new cobra.Command for the base command
func NewCommand() *cobra.Command {
	viper.SetEnvPrefix("tweetscollect")

	cmd := &cobra.Command{
		Use:   "tweetscollect -b <bearer-token> -t <topic> -o <output-file>",
		Short: "tweetscollect is a tool for collecting tweets for a topic.",
		Long:  `tweetscollect is a tool for collecting tweets for a topic from Twitter.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.CommandLine.Parse([]string{}); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract the arguments.
			bearerToken := viper.GetString(vBearerToken)
			outputFile := viper.GetString(vOutputFile)
			topic := viper.GetString(vTopic)
			if bearerToken == "" || topic == "" || outputFile == "" {
				return cmd.Usage()
			}

			return collectTweets(bearerToken, topic, outputFile)
		},
	}

	cmd.PersistentFlags().CountP("verbose", "v", "Increase verbosity. May be given multiple times.")

	cmd.Flags().StringP(flagBearerToken, shortHandBearerToken, "", "Bearer token")
	_ = viper.BindPFlag(vBearerToken, cmd.Flags().Lookup(flagBearerToken))

	cmd.Flags().StringP(flagOutputFile, shortHandOutputFile, "", "Output file")
	_ = viper.BindPFlag(vOutputFile, cmd.Flags().Lookup(flagOutputFile))

	cmd.Flags().StringP(flagTopic, shortHandTopic, "", "Topic, e.g. Facebook")
	_ = viper.BindPFlag(vTopic, cmd.Flags().Lookup(flagTopic))

	return cmd
}

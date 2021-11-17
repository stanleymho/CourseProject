package main

import (
	"context"
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
	vAccesssKeyID    = "access-key-id"
	vInputFile       = "input-file"
	vOutputFile      = "output-file"
	vLang            = "lang"
	vRegion          = "region"
	vSecretAccessKey = "secret-access-key"

	flagAccesssKeyID    = "access-key-id"
	flagInputFile       = "input-file"
	flagOutputFile      = "output-file"
	flagLang            = "lang"
	flagRegion          = "region"
	flagSecretAccessKey = "secret-access-key"

	shortHandAccesssKeyID    = "a"
	shortHandInputFile       = "i"
	shortHandOutputFile      = "o"
	shortHandLang            = "l"
	shortHandRegion          = "r"
	shortHandSecretAccessKey = "s"
)

// NewCommand returns a new cobra.Command for the base command
func NewCommand() *cobra.Command {
	viper.SetEnvPrefix("sentimentalyze")

	cmd := &cobra.Command{
		Use:   "sentimentalyze -i <input-file> -o <output-file> -a <access-key-id> -s <secret-access-key> [-r <region>] [-l <lang>]",
		Short: "sentimentalyze is a tool for performing sentiment analysis using Amazon Comprehend.",
		Long:  `sentimentalyze is a tool for performing sentiment analysis using Amazon Comprehend.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := flag.CommandLine.Parse([]string{}); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract the arguments.
			accessKeyID := viper.GetString(vAccesssKeyID)
			secretAccessKey := viper.GetString(vSecretAccessKey)
			inputFile := viper.GetString(vInputFile)
			outputFile := viper.GetString(vOutputFile)
			region := viper.GetString(vRegion)
			if inputFile == "" || outputFile == "" || accessKeyID == "" ||
				secretAccessKey == "" || region == "" {
				return cmd.Usage()
			}

			ctx := context.Background()
			return analyzeSentiment(ctx, inputFile, outputFile, region, accessKeyID, secretAccessKey)
		},
	}

	cmd.PersistentFlags().CountP("verbose", "v", "Increase verbosity. May be given multiple times.")

	cmd.Flags().StringP(flagAccesssKeyID, shortHandAccesssKeyID, "", "Access key ID")
	_ = viper.BindPFlag(vAccesssKeyID, cmd.Flags().Lookup(flagAccesssKeyID))

	cmd.Flags().StringP(flagInputFile, shortHandInputFile, "tweets.json", "Input file")
	_ = viper.BindPFlag(vInputFile, cmd.Flags().Lookup(flagInputFile))

	cmd.Flags().StringP(flagOutputFile, shortHandOutputFile, "tweets-sentiment.json", "Output file")
	_ = viper.BindPFlag(vOutputFile, cmd.Flags().Lookup(flagOutputFile))

	cmd.Flags().StringP(flagRegion, shortHandRegion, "us-east-1", "Region")
	_ = viper.BindPFlag(vRegion, cmd.Flags().Lookup(flagRegion))

	cmd.Flags().StringP(flagSecretAccessKey, shortHandSecretAccessKey, "", "Secret access key")
	_ = viper.BindPFlag(vSecretAccessKey, cmd.Flags().Lookup(flagSecretAccessKey))

	return cmd
}

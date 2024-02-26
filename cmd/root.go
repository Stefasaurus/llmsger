/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	process "github.com/Stefasaurus/llmsger/processing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CsvPath string
var OutputPath string
var MergeName string
var cfgFile string
var bMergeFiles bool

const mergeNameDefault string = "lang_merged"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "llmsger",
	Short: "llmsger is a tool for localization file generation.",
	Long: `
	llmsger CLI tool for automatic language localization file creation.
	It uses a provided CSV(.csv) file written by the user, to generate one merged 
	or multiple unmerged localization files for a programming language 
	(currently only supports C file generation)`,
	Example: `
Process a file so that only one language option can be active:
	.\llmsger.exe -f "sample3.csv" --mrg -n "stef"
Process a file so that all language option can be active:
	.\llmsger.exe dyngen -f "sample.csv" --basename "example" --varname "exampleVarname"`,
	Version: `0.2.0`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		defer func() {
			if err != nil {
				err = fmt.Errorf("llmsger failed: %w", err)
			}
		}()

		if CsvPath == "" {
			err = errors.New("path unspecified")
		} else {

			if cmd.Flags().Changed("mergename") && !bMergeFiles {
				err = fmt.Errorf("merged name cannot be used without the merge flasg (mrg)")
				return err
			}

			fmt.Println("Processing CSV...")

			//Process the CSV
			var langMap map[string][]string
			langMap, err = process.Csv(CsvPath)
			if err != nil {
				return err
			}

			if bMergeFiles {
				fmt.Println("Processing merged files...")
				//return nil
				err = process.CreateFilesMerge(langMap, OutputPath, MergeName)
			} else {
				err = process.CreateFiles(langMap, OutputPath)
			}

		}

		return err
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&CsvPath, "filepath", "f", "", "Input CSV file path")
	rootCmd.MarkPersistentFlagRequired("filepath")

	rootCmd.PersistentFlags().StringVarP(&OutputPath, "outdir", "o", ".", "Output file(s) path (optional)")

	rootCmd.Flags().StringVarP(&MergeName, "mergename", "n", mergeNameDefault, "Output file(s) path (optional)")

	rootCmd.Flags().BoolVar(&bMergeFiles, "mrg", false, "Output files will be merged into one file, by default the name is \"lang_merged\"")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".llmsger" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".llmsger")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

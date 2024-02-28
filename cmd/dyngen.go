/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	process "github.com/Stefasaurus/llmsger/processing"
	"github.com/spf13/cobra"
)

var DynBaseName string

var DynVariableName string

// dyngenCmd represents the dyngen command
var dyngenCmd = &cobra.Command{
	Use:   "dyngen",
	Short: "Generates dynamically changeable language files from CSV file.",
	Long: `
	Generates required files and text definitions from input CSV file.
	These files contain definitions which can be changed during runtime.
	All texts are contained within a variable that is included using the generated header file.
	For the C language generates a header and source file containing the main variable.`,
	Example: `.\llmsger.exe dyngen -f "sample.csv" --basename "example" --varname "exampleVarname"`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		defer func() {
			if err != nil {
				err = fmt.Errorf("dyngen failed: %w", err)
			}
		}()

		if CsvPath == "" {
			err = errors.New("path unspecified")
		} else {

			fmt.Println("Processing CSV...")

			//Process the CSV
			var langMap map[string][]string
			langMap, err = process.Csv(CsvPath)
			if err != nil {
				return err
			}

			if bReplace {
				err = process.ReplaceFields(langMap)
				if err != nil {
					return err
				}
			} else {
				delete(langMap, "from")
				delete(langMap, "toascii")
			}

			fmt.Println("Generating files...")
			//return nil
			err = process.CreateFilesDynamic(langMap, OutputPath, DynVariableName, DynBaseName)

		}
		fmt.Println("Done!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(dyngenCmd)
	dyngenCmd.MarkPersistentFlagRequired("filepath")

	dyngenCmd.Flags().StringVar(&DynBaseName, "basename", "llmsger", "Shared output file name (base name for source and header file) (optional)")

	dyngenCmd.Flags().StringVar(&DynVariableName, "varname", "gp_lang_texts", "Name of the variable for languages (optional)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dyngenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dyngenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

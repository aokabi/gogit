/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/aokabi/gogit/pkg"
	"github.com/spf13/cobra"
)

// updateIndexCmd represents the updateIndex command
var updateIndexCmd = &cobra.Command{
	Use:   "update-index",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if addFlag {
			if len(args) == 0 {
				fmt.Println("no file name")
				return
			}

			// create git object
			nameObjectMap := make(map[string]*pkg.GitObj)
			for _, arg := range args {
				obj := createObject(arg)
				nameObjectMap[arg] = obj
			}

			pkg.AddEntry(nameObjectMap)
		}

	},
}

var (
	addFlag bool
)

func init() {
	rootCmd.AddCommand(updateIndexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateIndexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateIndexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateIndexCmd.Flags().BoolVar(&addFlag, "add", false, "add file to index")
}

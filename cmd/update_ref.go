/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/aokabi/gogit/pkg"
	"github.com/spf13/cobra"
)

// updateRefCmd represents the updateRef command
var updateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("need 2 args")
			return
		}
		ref := args[0]
		newvalue := args[1]

		pkg.UpdateRefs(ref, newvalue)

	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateRefCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateRefCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

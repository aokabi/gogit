/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/aokabi/gogit/pkg"
	"github.com/aokabi/gogit/pkg/config"
	"github.com/spf13/cobra"
)

// commitTreeCmd represents the commitTree command
var commitTreeCmd = &cobra.Command{
	Use:   "commit-tree tree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.Read()
		now := time.Now()
		commit := pkg.NewCommit(
			args[0],
			conf.GetName(),
			conf.GetEmail(),
			now,
			conf.GetName(),
			conf.GetEmail(),
			now,
			commitFlag,
		)

		obj := commit.EncodeCommit()
		obj.Store()

		fmt.Println(obj.Hash())
	},
	Args: cobra.ExactArgs(1),
}

var (
	commitFlag string
)

func init() {
	rootCmd.AddCommand(commitTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitTreeCmd.Flags().StringVarP(&commitFlag, "m", "m", "", "commit message")
}

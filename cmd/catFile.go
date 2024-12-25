/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"compress/zlib"
	"fmt"
	"github.com/aokabi/gogit/pkg"
	"os"

	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		prefix := printFlag[:2]
		sufix := printFlag[2:]

		path := fmt.Sprintf(".git/objects/%s/%s", prefix, sufix)
		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// decompress the file
		r, err := zlib.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer r.Close()
		gitObj, err := pkg.Parse(r)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(gitObj.GetContent()))
	},
}

var (
	printFlag string
)

func init() {
	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	catFileCmd.PersistentFlags().StringVarP(&printFlag, "p", "p", "", "print file")
}

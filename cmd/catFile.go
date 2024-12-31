/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/aokabi/gogit/pkg"

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
		f := pkg.ReadObjectFile(printFlag)

		// decompress the file
		r, err := pkg.Decompress(f)
		if err != nil {
			panic(err)
		}
		defer r.Close()

		gitObj, err := pkg.Parse(r)
		if err != nil {
			panic(err)
		}
		fmt.Println(printString(gitObj))
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

func printString(o *pkg.GitObj) string {
	switch o.GetObjType() {
	case "blob":
		return o.DecodeContent2Blob()
	case "tree":
		tree := pkg.DecodeContent2Tree(o)
		entries := make([]string, 0)
		for e := range tree.Entries() {
			entries = append(entries, fmt.Sprintf("%s %s %s    %s", e.GetPerm(), e.GetObjType(), e.GetHash(), e.GetFilename()))
		}

		return strings.Join(entries, "\n")
	default:
		return fmt.Sprintf("unknown object %s", o.GetObjType())
	}
}

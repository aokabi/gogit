/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/aokabi/gogit/pkg"
	"github.com/spf13/cobra"
)

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object [-w] file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("no file name")
			return
		}

		createBlobObject(args[0])
	},
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.AddCommand(hashObjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	hashObjectCmd.Flags().BoolP("w", "w", false, "write database")
}

func createBlobObject(path string) *pkg.GitObj {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	obj := pkg.NewGitObj(pkg.BLOB, content)

	obj.Store()

	return obj
}

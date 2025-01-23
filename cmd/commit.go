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

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read index
		index, err := pkg.ReadIndexFile()
		if err != nil {
			fmt.Println(err)
			return
		}
		tree := pkg.NewTree()

		// indexのエントリをtreeのエントリに追加
		for entry := range index.Entries() {
			// 一旦全部blobとする
			tree.AddEntry(entry.GetPerm(), pkg.BLOB, entry.GetHash(), entry.GetFilename())
		}

		// create tree object
		treeObject := tree.EncodeTree()
		treeObject.Store()

		// get parent commit
		head := pkg.ReadHEAD()
		currentCommitHash := pkg.ReadRef(head)

		// create commit object
		conf := config.Read()
		now := time.Now()
		commit := pkg.NewCommit(
			treeObject.Hash(),
			currentCommitHash,
			conf.GetName(),
			conf.GetEmail(),
			now,
			conf.GetName(),
			conf.GetEmail(),
			now,
			messageFlag,
		)

		commitObject := commit.EncodeCommit()
		commitObject.Store()

		// update current ref
		pkg.UpdateRefs(head, commitObject.Hash())

		// print
		/*
					[master (root-commit) e188f8f] test
					 1 file changed, 1 insertion(+)
			 		 create mode 100644 1.txt
		*/
		fmt.Printf("[%s] %s\n", commitObject.Hash(), messageFlag)


	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringVarP(&messageFlag, "m", "m", "", "commit message")
}

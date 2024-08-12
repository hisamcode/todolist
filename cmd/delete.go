/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/hisamcode/todolist/cli"
	"github.com/hisamcode/todolist/internal/database/csvrepo"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete task",
	Long:  `delete task`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least 1 arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cli.StdErr(errors.New("argument must be number"))
			return
		}

		file, fileID, err := csvrepo.OpenFile("./internal/database/files/")
		if err != nil {
			cli.StdErr(err)
			return
		}
		defer file.Close()
		defer fileID.Close()

		repo := csvrepo.NewTaskModel(file, fileID)
		if err := repo.DeleteByID(id); err != nil {
			cli.StdErr(err)
			return
		}

		file.Seek(0, io.SeekStart)
		err = cli.RenderReplaceHeader(file, os.Stdout, tabwriter.TabIndent, []string{"ID", "Task", "Created At", "Done"})
		if err != nil {
			cli.StdErr(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

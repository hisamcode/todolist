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

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as done",
	Long:  `Mark a task as done`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least 1 arg")
		}

		_, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("argument must be number")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			cli.StdErr(err)
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
		task, err := repo.FindByID(id)
		if err != nil {
			cli.StdErr(err)
			return
		}

		task.IsComplete = true

		err = repo.UpdateByID(id, *task)
		if err != nil {
			cli.StdErr(err)
			return
		}

		file.Seek(0, io.SeekStart)
		cli.RenderReplaceHeader(file, os.Stdout, tabwriter.TabIndent, []string{"ID", "Task", "Created At", "Done"})
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

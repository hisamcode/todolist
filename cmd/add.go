/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/hisamcode/todolist/cli"
	"github.com/hisamcode/todolist/internal/data"
	"github.com/hisamcode/todolist/internal/database/csvrepo"
	"github.com/spf13/cobra"
)

var NameTask string
var DescriptionTask string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new todolist",
	Long:  `Add new todolist`,
	Run: func(cmd *cobra.Command, args []string) {
		file, fileID, err := csvrepo.OpenFile("./internal/database/files/")
		if err != nil {
			cli.StdErr(err)
			return
		}
		defer file.Close()
		defer fileID.Close()

		repo := csvrepo.NewTaskModel(file, fileID)
		task := data.Task{
			Description: NameTask,
			CreatedAt:   time.Now(),
			IsComplete:  false,
		}
		err = repo.Create(task)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}

		file.Seek(0, io.SeekStart)
		err = cli.RenderReplaceHeader(file, os.Stdout, tabwriter.TabIndent, []string{"ID", "Task", "Created At"})
		if err != nil {
			cli.StdErr(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&NameTask, "name", "", "Name of task")
	addCmd.MarkFlagRequired("name")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

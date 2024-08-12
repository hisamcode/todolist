/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hisamcode/todolist/cli"
	"github.com/hisamcode/todolist/internal/database/csvrepo"
	"github.com/spf13/cobra"
)

var All bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List todolist",
	Long:  `List todolist`,
	Run: func(cmd *cobra.Command, args []string) {
		file, fileID, err := csvrepo.OpenFile("./internal/database/files/")
		if err != nil {
			cli.StdErr(err)
			return
		}
		defer file.Close()
		defer fileID.Close()

		repo := csvrepo.NewTaskModel(file, fileID)
		records, err := repo.List()
		if err != nil {
			cli.StdErr(err)
			return
		}

		header := []string{"ID", "Task", "Created At"}
		csv := ""
		if All {
			header = append(header, "Done")
			csv, err = csvrepo.RecordsToStringCSV(records, true)
			if err != nil {
				cli.StdErr(err)
				return
			}
		} else {
			csv, err = csvrepo.RecordsToStringCSV(records, false)
			if err != nil {
				cli.StdErr(err)
				return
			}
		}

		r := strings.NewReader(csv)
		err = cli.RenderReplaceHeader(r, os.Stdout, tabwriter.TabIndent, header)
		if err != nil {
			cli.StdErr(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&All, "all", false, "--all")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

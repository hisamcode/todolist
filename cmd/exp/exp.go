package main

import (
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hisamcode/todolist/cli"
)

func main() {
	csv := `ID,Description,CreatedAt,IsComplete
1,My new task,2024-07-27T16:45:19-05:00,true
2,Finish this video,2024-07-27T16:45:26-05:00,true
3,Find a video editor,2024-07-27T16:45:31-05:00,false
4,Find a video editor,,false`
	s := strings.NewReader(csv)
	headers := []string{"ID", "Task", "Created at", "Done", "asd"}
	cli.RenderReplaceHeader(s, os.Stdout, tabwriter.TabIndent, headers)

}

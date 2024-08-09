package cli_test

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
	"testing/fstest"
	"text/tabwriter"

	"github.com/hisamcode/todolist/cli"
)

func dataCSV(t *testing.T) string {
	t.Helper()

	return `ID,Description,CreatedAt,IsComplete
1,My new task,2024-07-27T16:45:19-05:00,true
2,Finish this video,2024-07-27T16:45:26-05:00,true
3,Find a video editor,2024-07-27T16:45:31-05:00,false
4,Find a video editor,,false`
}

func fsCSV(t *testing.T) *fstest.MapFS {
	t.Helper()

	files := fstest.MapFS{
		"database/data.csv": {Data: []byte(dataCSV(t))},
	}

	return &files
}

func TestRenderReplaceHeaeder_BodyIsExpected(t *testing.T) {

}

func TestRenderReplaceHeader_TotalColumnMatchToTotalHeaderColumn(t *testing.T) {
	t.Parallel()

	// total colum header disamain yang ada di dataCSV
	testCases := []struct {
		name    string
		headers []string
	}{
		{name: "0 header", headers: []string{}},
		{name: "1 header", headers: []string{"id"}},
		{name: "2 header", headers: []string{"id", "task"}},
		{name: "3 header", headers: []string{"id", "task", "time"}},
		{name: "4 header", headers: []string{"id", "task", "time", "done"}},
		{name: "5 header, header lebih banyak dari body", headers: []string{"id", "task", "time", "done", "1"}},
		{name: "6 header, header lebih banyak dari body", headers: []string{"id", "task", "time", "done", "1", "2"}},
		{name: "7 header, header lebih banyak dari body", headers: []string{"id", "task", "time", "done", "1", "2", "3"}},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			files := fsCSV(t)
			file, err := files.Open("database/data.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			buf := new(bytes.Buffer)
			err = cli.RenderReplaceHeader(file, buf, tabwriter.Debug, v.headers)
			if err != nil {
				t.Fatal(err)
			}

			b := bufio.NewReader(buf)
			bline, _, err := b.ReadLine()
			if err != nil {
				t.Fatal(err)
			}
			sline := strings.Split(string(bline), "|")
			totalColumnHeader := len(sline)
			for {
				bline, _, err := b.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				sline := strings.Split(string(bline), "|")
				if totalColumnHeader != len(sline) {
					t.Errorf("want %d column but got %d column in body, string: %q", totalColumnHeader, len(sline), string(bline))
				}

			}
		})
	}

}

func TestRenderReplaceHeader_HeaderHaveNewLineAfterReplaceHeader(t *testing.T) {
	t.Parallel()

	files := fsCSV(t)
	file, err := files.Open("database/data.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	headers := []string{"ID", "Task"}
	err = cli.RenderReplaceHeader(file, buf, tabwriter.Debug, headers)
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 1)
	bb := []byte{}
	for {
		_, err := buf.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		bb = append(bb, b...)
		if string(b) == "\n" {
			break
		}
	}
	if !strings.Contains(string(bb), "Task\n") {
		t.Errorf("want have new line after header but got byte: %+v, string: %+v", bb, string(bb))
	}

}

func TestRenderReplaceHeader_HeaderOutputIsEqualToThePlan(t *testing.T) {
	t.Parallel()

	files := fsCSV(t)
	file, err := files.Open("database/data.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	headers := []string{"ID", "Task", "Created"}
	buf := new(bytes.Buffer)
	err = cli.RenderReplaceHeader(file, buf, tabwriter.Debug, headers)
	if err != nil {
		t.Fatal(err)
	}

	line, err := buf.ReadString('\n')
	if err != nil {
		return
	}

	wants := []string{"ID", "Task", "Created"}
	for _, v := range wants {
		if !strings.Contains(line, v) {
			t.Errorf("want %q but got %q", v, line)
		}
	}

}

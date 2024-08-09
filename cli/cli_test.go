package cli_test

import (
	"encoding/csv"
	"log"
	"testing"
	"testing/fstest"
)

func dataCSV(t *testing.T) string {
	t.Helper()

	dataCSV := `ID,Description,CreatedAt,IsComplete
	1,My new task,2024-07-27T16:45:19-05:00,true
	2,Finish this video,2024-07-27T16:45:26-05:00,true
	3,Find a video editor,2024-07-27T16:45:31-05:00,false`

	return dataCSV
}

func TestRender_HeaderOutputIsEqualToThePlan(t *testing.T) {
	t.Parallel()
	files := fstest.MapFS{
		"database/file.csv": {Data: []byte(dataCSV(t))},
	}

	file, err := files.Open("database/file.csv")
	if err != nil {
		t.Fatal(err)
	}

	cr := csv.NewReader(file)

	records, err := cr.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	want := "ID"
	got := records[0][0]
	if want != got {
		t.Errorf("want %q got %q", want, got)
	}

	want = "Description"
	got = records[0][1]
	if want != got {
		t.Errorf("want %q got %q", want, got)
	}

	want = "CreatedAt"
	got = records[0][2]
	if want != got {
		t.Errorf("want %q got %q", want, got)
	}

	want = "IsComplete"
	got = records[0][3]
	if want != got {
		t.Errorf("want %q got %q", want, got)
	}

}

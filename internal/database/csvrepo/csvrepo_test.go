package csvrepo_test

import (
	"encoding/csv"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/hisamcode/todolist/internal/data"
	"github.com/hisamcode/todolist/internal/database/csvrepo"
)

func tempCSVFile(t *testing.T) (*os.File, error) {
	t.Helper()

	return os.OpenFile(t.TempDir()+"csv.csv", os.O_CREATE|os.O_RDWR, fs.FileMode(os.O_RDWR))
}

func TestCreate(t *testing.T) {
	t.Parallel()

	// butuh open
	file, err := tempCSVFile(t)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	d := data.Task{ID: 0,
		Description: "this is new description",
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}

	repo := csvrepo.NewTaskModel(file)
	err = repo.Create(d)
	if err != nil {
		t.Fatal(err)
	}

	file.Seek(0, io.SeekStart)
	csvr := csv.NewReader(file)
	records, err := csvr.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	task, err := csvrepo.RecordCSVToStruct(records[0])
	if err != nil {
		t.Fatal(err)
	}

	want := d

	if want.ID != task.ID {
		t.Errorf("want id %d, but got id %d", want.ID, task.ID)
	}
	if want.Description != task.Description {
		t.Errorf("want Description %q, but got Description %q", want.Description, task.Description)
	}
	if want.CreatedAt.Format(time.RFC3339) != task.CreatedAt.Format(time.RFC3339) {
		t.Errorf("want CreatedAt %+v, but got CreatedAt %+v", want.CreatedAt, task.CreatedAt)
	}
	if want.IsComplete != task.IsComplete {
		t.Errorf("want IsComplete %+v, but got IsComplete %+v", want.IsComplete, task.IsComplete)
	}

}

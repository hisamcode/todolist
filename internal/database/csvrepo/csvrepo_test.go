package csvrepo_test

import (
	"encoding/csv"
	"io"
	"io/fs"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hisamcode/todolist/internal/data"
	"github.com/hisamcode/todolist/internal/database/csvrepo"
)

func tempCSVFile(t *testing.T) ([]*os.File, error) {
	t.Helper()

	file, err := os.OpenFile(t.TempDir()+"csv.csv", os.O_CREATE|os.O_RDWR, fs.FileMode(os.O_RDWR))
	if err != nil {
		return nil, err
	}

	fileID, err := os.OpenFile(t.TempDir()+"id.text", os.O_CREATE|os.O_RDWR, fs.FileMode(os.O_RDWR))
	if err != nil {
		return nil, err
	}

	return []*os.File{file, fileID}, nil

}

func TestCreate_RecordIDIsIncreament(t *testing.T) {
	t.Parallel()

	files, err := tempCSVFile(t)
	file := files[0]
	fileID := files[1]
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer fileID.Close()

	for i := 0; i < 5; i++ {
		d1 := data.Task{
			Description: "data" + strconv.Itoa(i),
			CreatedAt:   time.Now(),
			IsComplete:  false,
		}

		repo := csvrepo.NewTaskModel(file, fileID)
		err = repo.Create(d1)
		if err != nil {
			t.Fatal(err)
		}
	}

	file.Seek(0, io.SeekStart)
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		task, err := csvrepo.RecordCSVToStruct(records[i])
		if err != nil {
			t.Fatal(err)
		}

		want := i + 1
		if task.ID != want {
			t.Errorf("want id %d, got id %d", want, task.ID)
		}
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	files, err := tempCSVFile(t)
	file := files[0]
	fileID := files[1]

	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer fileID.Close()

	d := data.Task{ID: 1,
		Description: "this is new description",
		CreatedAt:   time.Now(),
		IsComplete:  false,
	}

	repo := csvrepo.NewTaskModel(file, fileID)
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

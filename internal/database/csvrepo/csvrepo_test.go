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

	t.Cleanup(func() {
		defer file.Close()
		defer fileID.Close()
	})

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

func TestList(t *testing.T) {
	t.Parallel()

	files, err := tempCSVFile(t)
	if err != nil {
		t.Fatal(err)
	}
	file := files[0]
	fileID := files[1]

	repo := csvrepo.NewTaskModel(file, fileID)

	tasks := []data.Task{}
	for i := 0; i < 5; i++ {
		task := data.Task{
			Description: strconv.Itoa(i),
			CreatedAt:   time.Now(),
			IsComplete:  false,
		}
		tasks = append(tasks, task)

		err := repo.Create(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	records, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		record, err := csvrepo.RecordCSVToStruct(records[i])
		if err != nil {
			t.Fatal(err)
		}

		if i+1 != record.ID {
			t.Errorf("want id %d got %d", i+1, record.ID)
		}
		if tasks[i].Description != record.Description {
			t.Errorf("want description %q got %q", tasks[i].Description, record.Description)
		}
		if tasks[i].CreatedAt.Format(time.RFC3339) != record.CreatedAt.Format(time.RFC3339) {
			t.Errorf("want Created at %v got %v", tasks[i].CreatedAt.Format(time.RFC3339), record.CreatedAt.Format(time.RFC3339))
		}
		if tasks[i].IsComplete != record.IsComplete {
			t.Errorf("want is complete %t got %t", tasks[i].IsComplete, record.IsComplete)
		}

	}

}

func TestDeleteByID(t *testing.T) {
	t.Parallel()

	files, err := tempCSVFile(t)
	if err != nil {
		t.Fatal(err)
	}

	file := files[0]
	fileID := files[1]

	repo := csvrepo.NewTaskModel(file, fileID)

	//
	for i := 0; i < 5; i++ {
		err := repo.Create(data.Task{
			Description: strconv.Itoa(i),
			CreatedAt:   time.Now(),
			IsComplete:  false,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	findID := 2
	err = repo.DeleteByID(findID)
	if err != nil {
		t.Fatal(err)
	}

	records, err := repo.List()
	if err != nil {
		t.Fatal(err)
	}

	found := 0
	for _, record := range records {
		if record[0] == strconv.Itoa(findID) {
			found++
		}
	}

	if found > 0 {
		t.Error("DeleteByID still not deleting")
	}

}

func TestUpdateByID(t *testing.T) {
	t.Parallel()

	files, err := tempCSVFile(t)
	if err != nil {
		t.Fatal(err)
	}

	repo := csvrepo.NewTaskModel(files[0], files[1])

	for i := 0; i < 5; i++ {
		repo.Create(data.Task{
			Description: strconv.Itoa(i),
			CreatedAt:   time.Now(),
			IsComplete:  false,
		})
	}

	updatedID := 2
	updatedData := data.Task{
		Description: "haha",
		IsComplete:  true,
	}
	err = repo.UpdateByID(updatedID, updatedData)
	if err != nil {
		t.Fatal(err)
	}

	files[0].Seek(0, io.SeekStart)
	r := csv.NewReader(files[0])
	records, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	foundRecord := []string{}
	for _, record := range records {
		if record[0] == strconv.Itoa(updatedID) {
			foundRecord = record
			break
		}
	}

	if updatedData.Description != foundRecord[1] {
		t.Errorf("want description %q but got %q", updatedData.Description, foundRecord[1])
	}

	gotIsComplete, err := strconv.ParseBool(foundRecord[3])
	if err != nil {
		t.Fatalf("error when parsing bool: %+v", err)
	}

	if updatedData.IsComplete != gotIsComplete {
		t.Errorf("want is complete %t but got %t", updatedData.IsComplete, gotIsComplete)
	}

}

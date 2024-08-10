package csvrepo

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/hisamcode/todolist/internal/data"
)

type TaskModel struct {
	file *os.File
}

func RecordCSVToStruct(record []string) (*data.Task, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, record[2])
	if err != nil {
		return nil, err
	}

	b, err := strconv.ParseBool(record[3])
	if err != nil {
		return nil, err
	}

	return &data.Task{
		ID:          id,
		Description: record[1],
		CreatedAt:   t,
		IsComplete:  b,
	}, nil
}

// NewTaskModel create TaskModel
// use pointer bisa aja file nya besar
func NewTaskModel(file *os.File) *TaskModel {
	return &TaskModel{
		file: file,
	}
}

// Create is create a new task, write to csv
func (m TaskModel) Create(task data.Task) error {
	cw := csv.NewWriter(m.file)
	cw.Write([]string{strconv.Itoa(task.ID), task.Description, task.CreatedAt.Format(time.RFC3339), strconv.FormatBool(task.IsComplete)})
	cw.Flush()
	return nil

}

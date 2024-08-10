package csvrepo

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/hisamcode/todolist/internal/data"
)

type TaskModel struct {
	file   *os.File
	fileID *os.File
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
func NewTaskModel(file *os.File, fileID *os.File) *TaskModel {
	return &TaskModel{
		file:   file,
		fileID: fileID,
	}
}

// Create is create a new task, write to csv
func (m TaskModel) Create(task data.Task) error {
	id := 0

	m.fileID.Seek(0, io.SeekStart)
	contentFileID, err := io.ReadAll(m.fileID)
	if err != nil {
		return err
	}

	if string(contentFileID) == "" {
		id = 1
	} else {
		id, err = strconv.Atoi(string(contentFileID))
		if err != nil {
			return err
		}
		id++
	}

	cw := csv.NewWriter(m.file)
	cw.Write([]string{strconv.Itoa(id), task.Description, task.CreatedAt.Format(time.RFC3339), strconv.FormatBool(task.IsComplete)})
	cw.Flush()

	// replace fileID to id
	if err := m.fileID.Truncate(0); err != nil {
		return fmt.Errorf("error when truncate: %+v", err)
	}

	m.fileID.Seek(0, io.SeekStart)
	_, err = m.fileID.WriteString(strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("error when WriteString(fileID): %+v", err)
	}

	return nil

}

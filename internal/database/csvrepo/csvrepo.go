package csvrepo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
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
	// TODO seek file
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

	m.file.Seek(0, io.SeekEnd)
	// m.file.Write([]byte("\n"))
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

func (m TaskModel) List() ([][]string, error) {
	// TODO seek file

	m.file.Seek(0, io.SeekStart)

	r := csv.NewReader(m.file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (m TaskModel) DeleteByID(id int) error {
	// todo seek file
	m.file.Seek(0, io.SeekStart)

	records, err := m.List()
	if err != nil {
		return err
	}

	recordsSaved := [][]string{}
	for _, record := range records {
		rid, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		if rid != id {
			recordsSaved = append(recordsSaved, record)
		}
	}

	m.file.Seek(0, io.SeekStart)
	w := csv.NewWriter(m.file)
	err = w.WriteAll(recordsSaved)
	if err != nil {
		return err
	}

	return nil
}

func (m TaskModel) FindByID(id int) (*data.Task, error) {

	m.file.Seek(0, io.SeekStart)
	records, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		task, err := RecordCSVToStruct(record)
		if err != nil {
			return nil, err
		}
		if task.ID == id {
			return task, nil
		}
	}

	return nil, errors.New("record not found")

}

func (m TaskModel) UpdateByID(id int, task data.Task) error {
	// todo file seek
	m.file.Seek(0, io.SeekStart)
	records, err := m.List()
	if err != nil {
		return err
	}

	for _, record := range records {
		if record[0] == strconv.Itoa(id) {
			if record[1] != task.Description {
				record[1] = task.Description
			}
			isComplete, err := strconv.ParseBool(record[3])
			if err != nil {
				return err
			}
			if isComplete != task.IsComplete {
				record[3] = strconv.FormatBool(task.IsComplete)
			}
		}
	}

	m.file.Seek(0, io.SeekStart)
	w := csv.NewWriter(m.file)
	err = w.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func OpenFile(basePath string) (file *os.File, fileID *os.File, err error) {
	file, err = os.OpenFile(basePath+"csv.csv", os.O_RDWR, fs.FileMode(os.O_RDWR))
	if err != nil {
		return nil, nil, err
	}
	fileID, err = os.OpenFile(basePath+"id.csv", os.O_RDWR, fs.FileMode(os.O_RDWR))
	if err != nil {
		return nil, nil, err
	}

	return file, fileID, nil
}

func RecordsToStringCSV(records [][]string, showUnComplete bool) (string, error) {

	csvString := ""
	for _, record := range records {
		task, err := RecordCSVToStruct(record)
		if err != nil {
			return "", nil
		}
		if showUnComplete {
			csvString += strings.Join(record, ",") + "\n"
		} else {
			if !task.IsComplete {
				csvString += strings.Join(record, ",") + "\n"
			}
		}
	}

	return csvString, nil
}

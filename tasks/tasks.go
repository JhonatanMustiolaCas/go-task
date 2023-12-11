package task

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Date     string `json:"date"`
}

func ListTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("There's no task yet")
		return
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "STATUS", "TASK", "DATE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, task := range tasks {

		status := "☐"
		if task.Complete == true {
			status = "☑"
		}
		tbl.AddRow(task.ID, status, task.Name, task.Date)
	}

	tbl.Print()
}

func DeleteTask(tasks []Task, id int) []Task {
	for i, task := range tasks {
		if task.ID == id {
			return append(tasks[:i], tasks[i+1:]...)
		}
	}
	return tasks
}

func CompleteTask(tasks []Task, id int) []Task {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Complete = !task.Complete
			break
		}
	}
	return tasks
}

func Addtask(tasks []Task, name string) []Task {
	date := time.Now()

	newTask := Task{
		ID:       GetNextID(tasks),
		Name:     name,
		Complete: false,
		Date:     date.Format("2006-01-02 15:04:05"),
	}

	return append(tasks, newTask)
}

func SaveTask(file *os.File, tasks []Task) {
	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	err = file.Truncate(0)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(file)
	_, err = writer.Write(bytes)
	if err != nil {
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}

func GetNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].ID + 1
}

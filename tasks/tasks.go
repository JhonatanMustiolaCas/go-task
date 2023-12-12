package task

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

const TASK_FILE_NAME string = "/tasks.json"

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Date     string `json:"date"`
}

func ListTasks(path string) {
	tasks := getTasks(path)
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
		if task.Complete {
			status = "☑"
		}

		tbl.AddRow(task.ID, status, task.Name, task.Date)
	}

	tbl.Print()
}

func DeleteTask(id int, path string) {
	tasks := getTasks(path)
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	saveTask(path, tasks)
}

func CompleteTask(id int, path string) {
	tasks := getTasks(path)
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Complete = !task.Complete
			break
		}
	}
	saveTask(path, tasks)
}

func Addtask(name string, path string) {
	date := time.Now()
	tasks := getTasks(path)
	newTask := Task{
		ID:       getNextID(tasks),
		Name:     name,
		Complete: false,
		Date:     date.Format("2006-01-02 15:04:05"),
	}

	tasksUpdated := append(tasks, newTask)
	saveTask(path, tasksUpdated)
}

func saveTask(path string, tasks []Task) {
	bytes, err := json.Marshal(tasks)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(path+TASK_FILE_NAME, os.O_RDWR|os.O_CREATE, 0666)
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

func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	return tasks[len(tasks)-1].ID + 1
}

func getTasks(path string) []Task {
	taskFile, err := os.OpenFile(path+TASK_FILE_NAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer taskFile.Close()

	var tasks []Task

	info, err := taskFile.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() != 0 {
		bytes, err := io.ReadAll(taskFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &tasks)
		if err != nil {
			panic(err)
		}
	} else {
		tasks = []Task{}
	}

	return tasks
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	task "github.com/JhonatanMustiolaCas/go-tasks/tasks"
	workspace "github.com/JhonatanMustiolaCas/go-tasks/workspaces"
)

const TASK_FILE_NAME string = "tasks.json"

func main() {
	workspace.Setup()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	var tasks []task.Task
	var taskFile os.File

	switch os.Args[1] {
	case "list":
		if len(os.Args) == 4 && os.Args[2] == "-wks" {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Printf("-wks argument must be a numeric charecter")
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			tasks, _ = filter(wks.Path)
			task.ListTasks(tasks)
			return
		}
		tasks, _ = filter("")
		task.ListTasks(tasks)
	case "add":
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("What's your task?")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if (len(os.Args) == 4) && (os.Args[2] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Printf("-wks argument must be a numeric charecter")
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			tasks, taskFile = filter(wks.Path)
			tasks = task.Addtask(tasks, name)
			task.SaveTask(&taskFile, tasks)
			return
		}

		tasks, taskFile = filter("")
		tasks = task.Addtask(tasks, name)
		task.SaveTask(&taskFile, tasks)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <task id>")
			return
		}
		if (len(os.Args) == 4) && (os.Args[2] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Printf("-wks argument must be a numeric charecter")
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			tasks, taskFile = filter(wks.Path)
			tasks = task.DeleteTask(tasks, wksId)
			task.SaveTask(&taskFile, tasks)
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("task id must be a numeric character")
			return
		}

		tasks, taskFile = filter("")
		tasks = task.DeleteTask(tasks, id)
		task.SaveTask(&taskFile, tasks)
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: complete <task id>")
			return
		}

		if (len(os.Args) == 4) && (os.Args[2] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Printf("-wks argument must be a numeric charecter")
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			tasks, taskFile = filter(wks.Path)
			tasks = task.CompleteTask(tasks, wksId)
			task.SaveTask(&taskFile, tasks)
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("task id must be a numeric character")
			return
		}

		tasks, taskFile = filter("")
		tasks = task.CompleteTask(tasks, id)
		task.SaveTask(&taskFile, tasks)

	case "create":
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if len(os.Args) == 3 {
			path = os.Args[2]
			_, err := os.ReadDir(path)
			if os.IsNotExist(err) {
				fmt.Printf("Directory %s doesn't exist\n", path)
			} else if err != nil {
				panic(err)
			}

		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("What's the name of your new workspace?")
		wksName, _ := reader.ReadString('\n')
		wksName = strings.TrimSpace(wksName)

		workspace.CreateWks(wksName, path)
		fmt.Printf("New workspace created: %s\n", wksName)

	case "wks":
		if os.Args[2] == "list" {
			workspace.ListWks()
		}

	default:
		printUsage()

	}
}

func printUsage() {
	fmt.Println("Usage: go-tasks [ list | add | delete <id> | complete <id> ]")
}

func openTaskFile(tasksPath string) os.File {
	taskFile, err := os.OpenFile(tasksPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	return *taskFile
}

func mkTaskList(taskFile os.File) []task.Task {
	var tasks []task.Task

	info, err := taskFile.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() != 0 {
		bytes, err := io.ReadAll(&taskFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &tasks)
		if err != nil {
			panic(err)
		}
	} else {
		tasks = []task.Task{}
	}
	return tasks
}

func filter(path string) ([]task.Task, os.File) {
	var err error
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	taskFile := openTaskFile(path + "/" + TASK_FILE_NAME)
	tasks := mkTaskList(taskFile)
	return tasks, taskFile
}

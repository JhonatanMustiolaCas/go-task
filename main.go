package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	task "github.com/jhonatanmustiolacas/go-tasks/tasks"
	workspace "github.com/jhonatanmustiolacas/go-tasks/workspaces"
)

const CWD = "."

func main() {
	workspace.Setup()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "list":
		if len(os.Args) == 4 && os.Args[2] == "-wks" {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Println("-wks argument must be a numeric character")
				return
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			task.ListTasks(wks.Path)
			return
		}
		task.ListTasks(CWD)
	case "add":
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("What's your task?")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		if (len(os.Args) == 4) && (os.Args[2] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Println("-wks argument must be a numeric charecter")
				return
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			task.Addtask(name, wks.Path)
			return
		}
		task.Addtask(name, CWD)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <task-id>")
			return
		}
		taskId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("task id must be a numeric charecter")
			return
		}
		if (len(os.Args) == 5) && (os.Args[3] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Println("-wks argument must be a numeric charecter")
				return
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			task.DeleteTask(taskId, wks.Path)
			return
		}
		task.DeleteTask(taskId, CWD)
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: complete <task id>")
			return
		}

		taskId, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("task id must be a numeric character")
			return
		}

		if (len(os.Args) == 5) && (os.Args[3] == "-wks") {
			wksId, err := strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Println("-wks argument must be a numeric charecter")
				return
			}
			wks, err := workspace.GetWks(wksId)
			if err != nil {
				panic(err)
			}
			task.CompleteTask(taskId, wks.Path)
			return
		}
		task.CompleteTask(taskId, CWD)

	case "create":
		path, err := os.Getwd()
		if err != nil {
			panic(err)
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
	fmt.Println("Usage: go-tasks [ wks | [list | add | delete <id> | complete <id>] ]")
}

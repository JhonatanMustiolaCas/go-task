package workspace

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

const WORKSPACE_DIR_NAME = "/.go-task"
const WORKSPACE_FILE_NAME = "/workspaces.json"

type WorkspaceError struct{}
type PathError struct{}

func (err *WorkspaceError) Error() string {
	return "Workspace doesn't exist"
}

func (err *PathError) Error() string {
	return "Path doesn't exist"
}

type Workspace struct {
	ID      int    `json:"id"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Created string `json:"created"`
	Edited  string `json:"edited"`
}

type Global struct {
	ID string `json:"id"`
}

func Setup() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(home+WORKSPACE_DIR_NAME, 0770)
	if os.IsExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
}

func CreateWks(name string, path string) {
	if exists := pathExists(path); !exists {
		panic(&WorkspaceError{})
	}
	created := time.Now()
	newWks := Workspace{
		ID:      GetNextID(),
		Path:    path,
		Name:    name,
		Created: created.Format("2006-01-02 15:04:05"),
		Edited:  "",
	}
	workspaces := getWorkspaces()
	workspaces = append(workspaces, newWks)
	saveWorkspaces(workspaces)
}

func DeleteWks(id int) {
	workspaces := getWorkspaces()

	if len(workspaces) == 0 {
		fmt.Println("There's no workspace to be deleted yet")
		return
	}
	for i, wks := range workspaces {
		if wks.ID == id {
			workspaces = append(workspaces[:i], workspaces[i+1:]...)
		}
	}
	saveWorkspaces(workspaces)
}

func GetWks(id int) (Workspace, error) {
	workspaces := getWorkspaces()
	for _, wks := range workspaces {
		if wks.ID == id {
			return wks, nil
		}
	}
	return Workspace{}, &WorkspaceError{}
}

func ListWks() {
	workspaces := getWorkspaces()

	if len(workspaces) == 0 {
		fmt.Println("There's no workspace yet")
		return
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "NAME", "PATH", "CREATED", "LAST EDITED")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, wks := range workspaces {
		tbl.AddRow(wks.ID, wks.Name, wks.Path, wks.Created, wks.Edited)
	}
	tbl.Print()
}

func GetNextID() int {
	workspaces := getWorkspaces()
	if len(workspaces) == 0 {
		return 1
	}
	return workspaces[len(workspaces)-1].ID + 1
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		panic(err)
	} else if os.IsNotExist(err) {
		return false
	}
	return true
}

func saveWorkspaces(workspaces []Workspace) {
	bytes, err := json.Marshal(workspaces)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(getWorkspaceDirectoryPath()+WORKSPACE_FILE_NAME, os.O_RDWR|os.O_CREATE, 0666)
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

func getWorkspaces() []Workspace {
	wksFile, err := os.OpenFile(getWorkspaceDirectoryPath()+WORKSPACE_FILE_NAME, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer wksFile.Close()

	var workspaces []Workspace

	info, err := wksFile.Stat()
	if err != nil {
		panic(err)
	}

	if info.Size() != 0 {
		bytes, err := io.ReadAll(wksFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &workspaces)
		if err != nil {
			panic(err)
		}
	} else {
		workspaces = []Workspace{}
	}
	return workspaces
}
func getWorkspaceDirectoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home + WORKSPACE_DIR_NAME

}

package Zwuiix

import (
	"bufio"
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Startup struct {
	Path string
}

func listDirectoryContents(dirPath string) ([]os.FileInfo, error) {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func DeleteFilePath(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	processList, err := ps.Processes()
	if err != nil {
		return err
	}

	for _, process := range processList {
		processExecutablePath := process.Executable()
		if processExecutablePath == absPath {
			targetPID := process.Pid()
			targetProcess, err := os.FindProcess(targetPID)
			if err != nil {
				return err
			}

			err = targetProcess.Kill()
			if err != nil {
				return err
			}

			err = os.RemoveAll(absPath)
			if err != nil {
				return err
			}

			return nil
		}
	}

	err = os.RemoveAll(absPath)
	if err != nil {
		return err
	}

	return nil
}

func (s Startup) HasInfected() {
	entries, err := listDirectoryContents(s.Path)
	if err != nil {
		log.Fatal(err)
	}

	var suspected []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() || entry.Size() > 50000 || !(strings.Contains(entry.Name(), ".lnk") || strings.Contains(entry.Name(), ".ini")) {
			suspected = append(suspected, entry)
		}
	}

	if len(suspected) == 0 {
		fmt.Println("Following our analysis, you do not seem to have been infected in the startup!")
		return
	}

	if len(suspected) != 0 {
		for i := 0; i < len(suspected); i++ {
			entry := suspected[i]
			if entry.Size() > 50000 || strings.Contains(entry.Name(), ".exe") {
				fmt.Println("The " + entry.Name() + " file is very suspicious")
			} else {
				fmt.Println("The " + entry.Name() + " file is suspicious")
			}
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Startup infected. Do you want to delete unwanted files? (y/n): ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			fmt.Println("Delete infected files.")
			for i := 0; i < len(suspected); i++ {
				path := filepath.Join(s.Path, suspected[i].Name())
				err := DeleteFilePath(path)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else if answer == "n" || answer == "N" {
			fmt.Println("Continue with infected files!")
		}
	}
}

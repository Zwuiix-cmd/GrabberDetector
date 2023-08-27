package Injection

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Discord struct{}

func searchFoldersForKeyword(rootPath string, keyword string) ([]string, error) {
	var result []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		if strings.Contains(info.Name(), keyword) {
			result = append(result, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func searchFoldersWithFiles(rootPath string, filenames ...string) ([]string, error) {
	var result []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		foundAllFiles := true
		for _, filename := range filenames {
			filePath := filepath.Join(path, filename)
			if _, err := os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					foundAllFiles = false
					break
				}
				return err
			}
		}

		if foundAllFiles {
			result = append(result, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func getAllFoldersInDirectory(rootPath string) ([]string, error) {
	var folders []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			folders = append(folders, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return folders, nil
}

func writeToFile(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func stringInSlice(target string, slice []string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func detect() []string {
	discordFiles := []string{"Discord", "DiscordCanary", "DiscordPTB", "DiscordDevelopment"}
	var invalidPayload []string
	for i := 0; i < len(discordFiles); i++ {
		path := filepath.Join(`C:\\Users\\` + os.Getenv("username") + `\\AppData\\Local\\` + discordFiles[i])
		localPath, err := os.Stat(path)
		if err == nil && localPath.IsDir() {
			directory, err := getAllFoldersInDirectory(path)
			if err == nil {
				for d := 0; d < len(directory); d++ {
					path := directory[d]
					files, err := searchFoldersWithFiles(path, "index.js", "package.json", "core.asar")
					if err == nil {
						for i := 0; i < len(files); i++ {
							fileTextPath := filepath.Join(files[i] + `\\index.js`)
							contentFile, err := readFile(fileTextPath)
							if err == nil && contentFile != "module.exports = require('./core.asar');" && contentFile != `module.exports = require("./core.asar");` {
								target := fileTextPath + "|==GRABBER-DETECT==|" + contentFile
								if !stringInSlice(target, invalidPayload) {
									invalidPayload = append(invalidPayload, target)
								}
							}
						}
					}
				}
			}
		}
	}

	return invalidPayload
}

func (receiver Discord) HasInjection() {
	hasDetection := detect()
	if len(hasDetection) == 0 {
		fmt.Println("You don't currently have an injection.")
		return
	}

	for i := 0; i < len(hasDetection); i++ {
		split := strings.Split(hasDetection[i], "|==GRABBER-DETECT==|")
		path := split[0]
		content := split[1]

		donerium, sure := Donerium{}.Flag(content)
		webhookLinks, links := OtherGrabber{}.Flag(content)
		if donerium != "" {
			if sure {
				fmt.Println("[STEALER] Donerium flagged, webhook: " + donerium + " (" + path + ")")
			} else {
				fmt.Println("[GRABBER] UnknownName flagged, webhook: " + donerium)
			}
		} else if len(links) != 0 {
			fmt.Println("[GRABBER] UnknownName, links found: \n" + strings.Join(links, ",\n"))
			if len(webhookLinks) != 0 {
				fmt.Println("\n[WEBHOOK] Founded webhook links: \n" + strings.Join(webhookLinks, ",\n"))
			} else {
				fmt.Println("\n[WEBHOOK] No webhook found!")
			}
		} else {
			fmt.Println("Incorrect payload in " + path)
		}
	}

	if len(hasDetection) != 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Injection(s) were found. Do you want to restore the default settings? (y/n): ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			fmt.Println("Delete infected files, and write the original ones.")
			writeDefaultFiles()
		} else if answer == "n" || answer == "N" {
			fmt.Println("Continue with infected files!")
		}
	}
}

func writeDefaultFiles() {
	discordFiles := []string{"Discord", "DiscordCanary", "DiscordPTB", "DiscordDevelopment"}
	for i := 0; i < len(discordFiles); i++ {
		path := filepath.Join(`C:\\Users\\` + os.Getenv("username") + `\\AppData\\Local\\` + discordFiles[i])
		localPath, err := os.Stat(path)
		if err == nil && localPath.IsDir() {
			directory, err := getAllFoldersInDirectory(path)
			if err == nil {
				for d := 0; d < len(directory); d++ {
					path := directory[d]
					files, err := searchFoldersWithFiles(path, "index.js", "package.json", "core.asar")
					if err == nil {
						for i := 0; i < len(files); i++ {
							fileTextPath := filepath.Join(files[i] + `\\index.js`)
							err := writeToFile(fileTextPath, "module.exports = require('./core.asar');")
							if err != nil {
								log.Println("Unable to repair infected file name in "+fileTextPath+" because ", err)
							}
						}
					}
				}
			}
		}
	}
}

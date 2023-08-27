package Zwuiix

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"gopkg.in/toast.v1"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type TempStealStorage struct {
	Icon embed.FS
}

var Delete = false

func (s TempStealStorage) Start(value bool) {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err == nil {
		addHostEntry("127.0.0.1", "api.epsilon1337.com")
		addHostEntry("127.0.0.1", "epsilon1337.com")
	}

	if value {
		Delete = true
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Automatically deleted when a file is suspicious? (y/n): ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			Delete = true
		}
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Lunch(filepath.Join("C:\\Windows\\Temp"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\AppData\\Local\\Temp"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\AppData\\Roaming"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\Desktop"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\Documents"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\Downloads"))
			s.Lunch(filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\Videos"))
		}
	}
}

func (s TempStealStorage) Lunch(path string) {
	contents, err := listDirectoryContents(path)
	if err != nil {
		return
	}

	var suspected []File
	for i := 0; i < len(contents); i++ {
		file := contents[i]

		hostname, _ := os.Hostname()
		if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(hostname)) || strings.Contains(strings.ToLower(file.Name()), "discord_backup_codes") || strings.Contains(strings.ToLower(file.Name()), "epsilon-") || file.Name() == ("epsilon-"+os.Getenv("username")) {
			suspected = append(suspected, File{Path: filepath.Join(path + "\\" + file.Name()), FileInfo: file})
			continue
		}

		if !file.IsDir() {
			continue
		}

		suspected = append(suspected, foundFile(filepath.Join(path+"\\"+file.Name()))...)
	}

	for i := 0; i < len(suspected); i++ {
		file := suspected[i]

		notification := toast.Notification{
			AppID:   "MusuiDectector",
			Title:   "An unwanted file has been found!",
			Message: file.FileInfo.Name() + " folder!",
			Actions: []toast.Action{
				{"protocol", "Protect me now", ""},
				{"protocol", "Well done!", ""},
			},
		}
		_ = notification.Push()

		if Delete {
			fmt.Println("Removing " + file.FileInfo.Name() + " folder, found in " + file.Path)
			_ = DeleteFilePath(file.GetPath())
		} else {
			fmt.Println(file.FileInfo.Name() + " folder, found in " + file.Path)
		}
	}
}

func (s TempStealStorage) FoundAndDeleteExe(path string) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error recovering processes:", err)
		return
	}

	for _, p := range processes {
		exePath, err := p.Exe()
		if err == nil {
			openFiles, err := p.OpenFiles()
			if err == nil {
				for _, file := range openFiles {
					if strings.Contains(file.Path, path) {
						_ = DeleteFilePath(exePath)
					}
				}
			}
		}
	}
}

func foundFile(path string) []File {
	contents, err := listDirectoryContents(path)
	if err != nil {
		return []File{}
	}

	keywords := []string{
		/* DONERIUM */
		"Autofill",
		"Bookmark",
		"Cookie",
		"History",
		"Password",
		"Steam Account",
		"TikTok Account",
		"Minecraft Account",
		"Growtopia",
		"Telegram",
		"Wallet",
		"Executable Info.txt",
		"Found Wallets.txt",
		"Network Data.txt",
		"User Info.txt",
		"WiFi Connections.txt",

		/* EPSILON */
		"Autofill Data",
		"Credit Card",
		"Credit",
		"Messengers",
	}
	var suspected []File
	for i := 0; i < len(contents); i++ {
		file := contents[i]

		hostname, _ := os.Hostname()
		if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(hostname)) || strings.Contains(strings.ToLower(file.Name()), "discord_backup_codes") || strings.Contains(strings.ToLower(file.Name()), "epsilon-") || file.Name() == ("epsilon-"+os.Getenv("username")) {
			suspected = append(suspected, File{Path: filepath.Join(path + "\\" + file.Name()), FileInfo: file})
			continue
		}

		if !file.IsDir() {
			continue
		}

		for i := 0; i < len(keywords); i++ {
			if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(keywords[i])) {
				suspected = append(suspected, File{Path: filepath.Join(path + "\\" + file.Name()), FileInfo: file})
				continue
			}
		}

		//suspected = append(suspected, foundFile(filepath.Join(path+"\\"+file.Name()))...)
	}
	return suspected
}

func addHostEntry(ip, domain string) {
	f, err := os.OpenFile("C:\\Windows\\System32\\drivers\\etc\\hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	entry := fmt.Sprintf("%s %s\n", ip, domain)
	_, err = f.WriteString(entry)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

package main

import (
	"GrabberDectector/src/Zwuiix"
	"GrabberDectector/src/Zwuiix/Injection"
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed resources/logo.png
var iconData embed.FS

func main() {
	Zwuiix.SetConsoleTitle("Grabber Detector")
	if runtime.GOOS != "windows" {
		fmt.Println("Sorry, this software is only functional on Windows.")
		Zwuiix.WaitForInterrupt()
		return
	}

	args := strings.Join(os.Args, ", ")
	if strings.Contains(args, "-startup") {
		Zwuiix.TempStealStorage{Icon: iconData}.Start(true)
		return
	}

	fmt.Println("Welcome to Grabber Detector, by @Zwuiix-cmd on Github")
	fmt.Println()
	fmt.Println("Explication: We currently have two options, one will be to analyze files from your computer, the other will analyze all the packets that go to discord!")
	fmt.Println()
	fmt.Println(" - Choice 1 => This option will look if you have infected files in your discord or your startup, and a deletion will be performed (optional)")
	fmt.Println(" - Choice 2 => This option will launch a scan of all the packets you send to discord, close all browsers on your computer and then launch the scan, if requests are made when you change your password, launch software, launch discord (software), this means that the application makes requests to discord, probably webhook requests.")
	fmt.Println(" - Choice 3 => Detect all storage of passwords, cookies and... in the pc's time, this prevents steelers from compressing their zip with the information!")
	fmt.Println(" - Choice 4 => Choice 3 And Install AutoStartup")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nChoice (1/2/3/4): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	if answer == "1" {
		fmt.Println("\nDiscord Infection Check!")
		Injection.Discord{}.HasInjection()

		path := filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup")
		fmt.Println("\nStartup Infection Check! (" + path + ")")
		startup := Zwuiix.Startup{Path: path}
		startup.HasInfected()
		Zwuiix.WaitForInterrupt()
		return
	} else if answer == "2" {
		Zwuiix.PacketScanner{}.Start()
	} else if answer == "3" {
		Zwuiix.TempStealStorage{Icon: iconData}.Start(false)
	} else if answer == "4" {
		startupFolder := filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup")
		desktopFolder := filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\Desktop")
		startMenuFolder := filepath.Join("C:\\Users\\" + os.Getenv("username") + "\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs")
		path := filepath.Join("C:\\Program Files\\Grabber Detector")

		currentExePath, err := os.Executable()
		if err != nil {
			fmt.Println("Error retrieving executable path:", err)
			Zwuiix.WaitForInterrupt()
			return
		}

		if filepath.Dir(currentExePath) == path {
			fmt.Println("Sorry, you have already performed an installation!")
			Zwuiix.WaitForInterrupt()
			return
		}

		_, err = os.Open("\\\\.\\PHYSICALDRIVE0")
		if err != nil {
			fmt.Println("I need additional permissions to access the " + path + "! (Reopen me as Administrator)")
			Zwuiix.WaitForInterrupt()
			return
		}

		if err := Zwuiix.CreateFolderIfNotExists(path); err != nil {
			fmt.Println("Error creating folder:", err)
			Zwuiix.WaitForInterrupt()
			return
		}

		newExePath := filepath.Join(path, filepath.Base(currentExePath))
		if err := os.Rename(currentExePath, newExePath); err != nil {
			fmt.Println("Error moving executable:", err)
			Zwuiix.WaitForInterrupt()
			return
		}
		linkPath := filepath.Join(startupFolder, "Grabber Detector.lnk")
		if err := Zwuiix.CreateShortcut(linkPath, newExePath, "Grabber and Stealer Detector", "-startup"); err != nil {
			fmt.Println("Error creating shortcut:", err)
			Zwuiix.WaitForInterrupt()
			return
		}

		desktopPath := filepath.Join(desktopFolder, "Grabber Detector.lnk")
		if err := Zwuiix.CreateShortcut(desktopPath, newExePath, "Grabber and Stealer Detector"); err != nil {
			fmt.Println("Error creating shortcut:", err)
			Zwuiix.WaitForInterrupt()
			return
		}

		startMenuPath := filepath.Join(startMenuFolder, "Grabber Detector.lnk")
		if err := Zwuiix.CreateShortcut(startMenuPath, newExePath, "Grabber and Stealer Detector"); err != nil {
			fmt.Println("Error creating shortcut:", err)
			Zwuiix.WaitForInterrupt()
			return
		}

		Zwuiix.TempStealStorage{Icon: iconData}.Start(true)
	}
}

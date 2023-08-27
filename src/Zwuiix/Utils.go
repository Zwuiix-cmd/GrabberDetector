package Zwuiix

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"
)

const (
	SW_HIDE = 0
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procShowWindow = modkernel32.NewProc("ShowWindow")
)

func CreateFolderIfNotExists(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(folderPath, os.ModePerm)
	}
	return err
}

func CreateShortcut(linkPath, targetPath string, description string, flags ...string) error {
	_ = ole.CoInitialize(0)
	defer ole.CoUninitialize()

	clsid, err := ole.CLSIDFromProgID("WScript.Shell")
	if err != nil {
		return err
	}

	unknown, err := ole.CreateInstance(clsid, nil)
	if err != nil {
		return err
	}
	defer unknown.Release()

	wshell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", linkPath)
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	defer idispatch.Release()

	_, err = oleutil.PutProperty(idispatch, "TargetPath", targetPath)
	if err != nil {
		return err
	}

	args := fmt.Sprintf(`"%s" `+strings.Join(flags, " "), targetPath)
	_, err = oleutil.PutProperty(idispatch, "Arguments", args)
	if err != nil {
		return err
	}

	_, err = oleutil.PutProperty(idispatch, "Description", description)
	if err != nil {
		return err
	}

	_, err = oleutil.CallMethod(idispatch, "Save")
	if err != nil {
		return err
	}

	return nil
}

func HideConsoleWindow() {
	_, _, _ = procShowWindow.Call(os.Stdout.Fd(), uintptr(SW_HIDE))
}

func WaitForInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
	<-interrupt
}

func SetConsoleTitle(title string) {
	kernel32, _ := syscall.LoadDLL("kernel32.dll")
	setConsoleTitle, _ := kernel32.FindProc("SetConsoleTitleW")

	utf16Title := syscall.StringToUTF16Ptr(title)
	_, _, _ = setConsoleTitle.Call(uintptr(unsafe.Pointer(utf16Title)))
}

package Zwuiix

import "os"

type File struct {
	Path     string
	FileInfo os.FileInfo
}

func (f File) GetPath() string {
	return f.Path
}

func (f File) GetFileInfo() os.FileInfo {
	return f.FileInfo
}

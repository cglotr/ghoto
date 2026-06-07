package fileremover

import "os"

type File_remover__impl struct{}

func File_remover__impl__new() *File_remover__impl {
	return &File_remover__impl{}
}

func (f *File_remover__impl) Remove(file string) error {
	return os.Remove(file)
}

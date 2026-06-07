package fileremover

import "fmt"

type File_remover__dummy struct{}

func File_remover__dummy__new() *File_remover__dummy {
	return &File_remover__dummy{}
}

func (f *File_remover__dummy) Remove(file string) error {
	fmt.Printf("🤡 File_remover__dummy__Remove: %v\n", file)
	return nil
}

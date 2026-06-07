package fileremover

type File_remover interface {
	Remove(file string) error
}

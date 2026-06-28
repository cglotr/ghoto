package ghoto

import "path"

type Photo_upload struct {
	Order        int
	File_path    string
	Upload_token string
}

func (p *Photo_upload) Get_file_name() string {
	return path.Base(p.File_path)
}

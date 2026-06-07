package googlephotos

import "errors"

type Google_photos__dummy struct {
}

func Google_photos__dummy__new() *Google_photos__dummy {
	return &Google_photos__dummy{}
}

func (g *Google_photos__dummy) Create_album(album_name string) (*Google_album, error) {
	return &Google_album{
		Id:   "google_album_id",
		Name: "google_album_name",
	}, nil
}

func (g *Google_photos__dummy) Get_photo(photo_id string) (*Google_photo, error) {
	return nil, errors.New("Get_photo")
}

func (g *Google_photos__dummy) List_album() (*Res__list_album, error) {
	return &Res__list_album{
		Albums: []Res__album{
			{
				Id:    "album_id",
				Title: "Insta360",
			},
		},
	}, nil
}

func (g *Google_photos__dummy) Upload_photo(file_path string, google_album Google_album) (*Google_photo, error) {
	return &Google_photo{
		Id:       "id__" + file_path,
		Filename: "file_name__" + file_path,
	}, nil
}

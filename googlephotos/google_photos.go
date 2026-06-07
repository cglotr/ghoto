package googlephotos

type Google_photos interface {
	Create_album(album_name string) (*Google_album, error)
	Get_photo(photo_id string) (*Google_photo, error)
	List_album() (*Res__list_album, error)
	Upload_photo(file_path string, google_album Google_album) (*Google_photo, error)
}

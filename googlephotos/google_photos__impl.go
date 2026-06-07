package googlephotos

import "net/http"

type Google_photos__impl struct {
	client *http.Client
}

func Google_photos__impl__new(client *http.Client) *Google_photos__impl {
	return &Google_photos__impl{
		client: client,
	}
}

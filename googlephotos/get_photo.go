package googlephotos

import (
	"encoding/json"
	"io"
)

func (g *Google_photos__impl) Get_photo(photo_id string) (*Google_photo, error) {
	var err error

	url := "https://photoslibrary.googleapis.com/v1/mediaItems/" + photo_id
	res, err := g.client.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var res__mediaItem Res__mediaItem
	err = json.Unmarshal(b, &res__mediaItem)
	if err != nil {
		return nil, err
	}

	return &Google_photo{
		Id:         res__mediaItem.Id,
		ProductUrl: res__mediaItem.ProductUrl,
		Filename:   res__mediaItem.Filename,
	}, nil
}

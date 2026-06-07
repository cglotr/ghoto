package googlephotos

import (
	"bytes"
	"encoding/json"
	"io"
)

func (g *Google_photos__impl) Create_album(album_name string) (*Google_album, error) {
	type Req__album struct {
		Title string `json:"title"`
	}

	type Req struct {
		Album Req__album `json:"album"`
	}

	req := Req{
		Album: Req__album{
			Title: album_name,
		},
	}

	url := "https://photoslibrary.googleapis.com/v1/albums"
	contentType := "application/json"
	b, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	res, err := g.client.Post(url, contentType, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type Res struct {
		Album Res__album `json:"album"`
	}

	var res_res Res
	err = json.Unmarshal(b, &res_res)
	if err != nil {
		return nil, err
	}

	return &Google_album{
		Id:   res_res.Album.Id,
		Name: res_res.Album.Title,
	}, nil
}

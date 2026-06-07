package googlephotos

import (
	"encoding/json"
	"io"
	"strconv"
)

func (g *Google_photos__impl) List_album() (*Res__list_album, error) {
	if g.client == nil {
		panic("List_album__client__nil")
	}

	res__album := []Res__album{}
	res__list_album, err := g.list_album("")
	if err != nil {
		return nil, err
	}
	res__album = append(res__album, res__list_album.Albums...)

	next_page_token := res__list_album.NextPageToken
	for next_page_token != "" {
		res__list_album, err = g.list_album(next_page_token)
		if err != nil {
			return nil, err
		}
		res__album = append(res__album, res__list_album.Albums...)
		next_page_token = res__list_album.NextPageToken
	}

	return &Res__list_album{
		Albums:        res__album,
		NextPageToken: next_page_token,
	}, nil
}

func (g *Google_photos__impl) list_album(pageToken string) (*Res__list_album, error) {
	if g.client == nil {
		panic("list_album__client__nil")
	}
	pageSize := 10
	url := "https://photoslibrary.googleapis.com/v1/albums"
	url += "?pageSize=" + strconv.Itoa(pageSize)
	if pageToken != "" {
		url += "&pageToken=" + pageToken
	}

	res, err := g.client.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res__list_album := Res__list_album{}
	err = json.Unmarshal(b, &res__list_album)
	if err != nil {
		return nil, err
	}

	return &res__list_album, nil
}

type Res__album struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	IsWriteable bool   `json:"isWriteable"`
}

type Res__list_album struct {
	Albums        []Res__album `json:"albums"`
	NextPageToken string       `json:"nextPageToken"`
}

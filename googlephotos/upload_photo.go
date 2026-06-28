package googlephotos

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cglotr/ghoto/constant"
)

func (g *Google_photos__impl) Create_photo(
	upload_token string,
	google_album Google_album,
	file_name string,
) (*Google_photo, error) {
	try_count := 0
	func__post := func() (*Res__mediaItem, error) {
		return g.create_photo(upload_token, google_album, file_name)
	}

	res__mediaItem, err := func__post()
	for err != nil && try_count < constant.Retry__count {
		try_count += 1
		time.Sleep(time.Duration(try_count) * 10 * time.Second)
		res__mediaItem, err = func__post()
	}

	if err != nil {
		return nil, errors.New("Upload_photo__create_photo__err: " + err.Error())
	}

	return &Google_photo{
		Id:         res__mediaItem.Id,
		ProductUrl: res__mediaItem.ProductUrl,
		Filename:   res__mediaItem.Filename,
	}, nil
}

func (g *Google_photos__impl) Upload_photo(file_path string) (*string, error) {
	var err error

	url := "https://photoslibrary.googleapis.com/v1/uploads"
	contentType := "application/octet-stream"
	b, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	try_count := 0
	func__post := func() (*http.Response, error) {
		return g.client.Post(url, contentType, bytes.NewBuffer(b))
	}

	res, err := func__post()
	for err != nil && try_count < constant.Retry__count {
		try_count += 1
		time.Sleep(time.Duration(try_count) * 10 * time.Second)
		res, err = func__post()
	}

	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	upload_token := string(b)

	return &upload_token, nil
}

func (g *Google_photos__impl) create_photo(
	upload_token string,
	google_album Google_album,
	file_name string,
) (*Res__mediaItem, error) {
	type Req__simpleMediaItem struct {
		UploadToken string `json:"uploadToken"`
		FileName    string `json:"fileName"`
	}
	type Req__newMediaItems struct {
		SimpleMediaItem Req__simpleMediaItem `json:"simpleMediaItem"`
	}
	type Req__albumPosition struct {
		Position string `json:"position"`
	}
	type Req__create_media_item struct {
		AlbumId       string               `json:"albumId"`
		NewMediaItems []Req__newMediaItems `json:"newMediaItems"`
		AlbumPosition Req__albumPosition   `json:"albumPosition"`
	}

	req := Req__create_media_item{
		AlbumId: google_album.Id,
		NewMediaItems: []Req__newMediaItems{
			{
				SimpleMediaItem: Req__simpleMediaItem{
					UploadToken: upload_token,
					FileName:    file_name,
				},
			},
		},
		AlbumPosition: Req__albumPosition{
			Position: "FIRST_IN_ALBUM",
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := g.client.Post(
		"https://photoslibrary.googleapis.com/v1/mediaItems:batchCreate",
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}

	type Res__status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	type Res__newMediaItemResults struct {
		Status    Res__status    `json:"status"`
		MediaItem Res__mediaItem `json:"mediaItem"`
	}
	type Res__media_items struct {
		NewMediaItemResults []Res__newMediaItemResults `json:"newMediaItemResults"`
	}

	var new_media_item_results Res__media_items
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &new_media_item_results)
	if err != nil {
		return nil, err
	}

	if len(new_media_item_results.NewMediaItemResults) < 1 {
		return nil, errors.New("Upload_photo__NewMediaItemResults__not_at_least__1")
	}
	media_item := new_media_item_results.NewMediaItemResults[0]

	if media_item.Status.Code != 0 {
		panic("create_photo__media_item_status_code__not_ok: " + strconv.Itoa(media_item.Status.Code))
	}

	return &media_item.MediaItem, nil
}

type Res__mediaItem struct {
	Id         string `json:"id"`
	ProductUrl string `json:"productUrl"`
	Filename   string `json:"filename"`
}

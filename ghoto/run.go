package ghoto

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/cglotr/ghoto/fileremover"
	"github.com/cglotr/ghoto/googleauth"
	"github.com/cglotr/ghoto/googlephotos"
	"github.com/cglotr/ghoto/util"
)

const ghoto_version = "v0"

type Ghoto struct {
	google_auth   googleauth.Google_auth
	google_photos googlephotos.Google_photos
	file_remover  fileremover.File_remover
}

func Ghoto__new() *Ghoto {
	return &Ghoto{
		google_auth:   googleauth.Google_auth__dummy__new(),
		google_photos: googlephotos.Google_photos__dummy__new(),
		file_remover:  fileremover.File_remover__dummy__new(),
	}
}

func (g *Ghoto) Activate() {
	g.google_auth = googleauth.Google_auth__impl__new()

	client, err := g.google_auth.Get_client()
	if err != nil {
		panic("Activate__Get_client:\n\t" + err.Error())
	}

	g.google_photos = googlephotos.Google_photos__impl__new(client)
	g.file_remover = fileremover.File_remover__impl__new()
}

func (g *Ghoto) Run(dir string, album_name string) error {
	dir = filepath.Dir(dir + "/")
	fmt.Printf("🌿 Ghoto %v: dir=%v, album=%v\n",
		ghoto_version,
		dir,
		album_name,
	)

	var google_album *googlephotos.Google_album
	res__list_album, err := g.google_photos.List_album()
	if err != nil {
		panic("Could not get album list: " + err.Error())
	}
	for _, album := range res__list_album.Albums {
		if album.Title == album_name {
			google_album = &googlephotos.Google_album{
				Id:   album.Id,
				Name: album.Title,
			}
		}
	}
	if google_album == nil {
		google_album, err = g.google_photos.Create_album(album_name)
		if err != nil {
			panic("Could not create Google Photos album: " + err.Error())
		}
	}

	sorted_photo_files := util.Sort_files(util.Filter_photo_files(util.Get_files(dir)))

	batch_cursor := 0
	for batch_cursor < len(sorted_photo_files) {
		photos_per_batch := 10
		batch_end := min(
			batch_cursor+photos_per_batch,
			len(sorted_photo_files),
		)
		batch_photo_files := sorted_photo_files[batch_cursor:batch_end]

		err = g.upload_photo_files(batch_photo_files, *google_album)
		if err != nil {
			return err
		}

		batch_cursor += photos_per_batch
	}

	photo_files := util.Filter_photo_files(util.Get_files(dir))
	if len(photo_files) > 0 {
		return errors.New("Photo files remaining!")
	}

	non_photo_files := util.Filter_non_photo_files(util.Get_files(dir))
	for _, non_photo_file := range non_photo_files {
		os.Remove(non_photo_file)
	}

	return nil
}

func (g *Ghoto) upload_photo_files(photo_files []string, google_album googlephotos.Google_album) error {
	wg := &sync.WaitGroup{}
	ch__photo_upload := make(chan Photo_upload, len(photo_files))

	for order, photo_file := range photo_files {
		file_for_worker := File_for_worker{
			Order:     order,
			File_path: photo_file,
		}

		wg.Add(1)
		go g.work__upload_photo(
			wg,
			file_for_worker,
			ch__photo_upload,
		)
	}

	wg.Wait()

	photo_uploads := []Photo_upload{}
	for range photo_files {
		photo_uploads = append(photo_uploads, <-ch__photo_upload)
	}
	slices.SortFunc(photo_uploads, func(a, b Photo_upload) int {
		return cmp.Compare(a.Order, b.Order)
	})

	check__failure := false
	for _, photo_upload := range photo_uploads {
		check__photo_upload := false
		check__photo_create := false
		check__photo_cleanup := false

		if photo_upload.Upload_token != "" {
			check__photo_upload = true
		}

		google_photo, err := g.google_photos.Create_photo(photo_upload.Upload_token, google_album)
		if err == nil && google_photo.ProductUrl != "" {
			check__photo_create = true
		}

		if check__photo_upload && check__photo_create {
			err = g.file_remover.Remove(photo_upload.File_path)
			if err == nil {
				check__photo_cleanup = true
			}
		}

		if check__photo_upload && check__photo_create && check__photo_cleanup {
			fmt.Printf("✅ Photo upload done: file=%v, url=%v\n",
				google_photo.Filename,
				google_photo.ProductUrl,
			)
		} else {
			check__failure = true
			fmt.Printf("❌ Photo upload failed: file=%v\n",
				photo_upload.File_path,
			)
		}
	}

	if check__failure {
		return errors.New("Some photos failed to upload")
	} else {
		return nil
	}
}

func (g *Ghoto) work__upload_photo(
	wg *sync.WaitGroup,
	file_for_worker File_for_worker,
	ch__photo_upload chan Photo_upload,
) {
	defer wg.Done()

	upload_token, err := g.google_photos.Upload_photo(file_for_worker.File_path)
	if err != nil {
		fmt.Printf("❌ Photo upload failed: file=%v\n",
			file_for_worker,
		)
		return
	}
	ch__photo_upload <- Photo_upload{
		Order:        file_for_worker.Order,
		File_path:    file_for_worker.File_path,
		Upload_token: *upload_token,
	}
}

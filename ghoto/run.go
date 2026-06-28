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
	dir = filepath.Dir(dir)
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

		err = g.run(batch_photo_files, album_name)
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

func (g *Ghoto) run(files []string, album_name string) error {
	worker_count := max(1, min(10, len(files)))
	files_per_worker := (len(files) / worker_count) + 1

	work_assigned_count := 0
	wg := &sync.WaitGroup{}
	ch__photo_upload := make(chan Photo_upload, len(files))
	photo_order := 0

	for worker_id := range worker_count {
		i := files_per_worker * worker_id
		if i >= len(files) {
			continue
		}
		j := min(i+files_per_worker, len(files))
		files_for_worker := []File_for_worker{}
		for _, file := range files[i:j] {
			file_for_worker := File_for_worker{
				Order:     photo_order,
				File_path: file,
			}
			photo_order += 1
			files_for_worker = append(files_for_worker, file_for_worker)
		}

		wg.Add(1)
		go g.work__upload_photo(
			wg,
			files_for_worker,
			ch__photo_upload,
		)

		work_assigned_count += j - i
		if work_assigned_count >= len(files) {
			break
		}
	}

	wg.Wait()

	photo_uploads := []Photo_upload{}
	for range files {
		photo_uploads = append(photo_uploads, <-ch__photo_upload)
	}
	slices.SortFunc(photo_uploads, func(a, b Photo_upload) int {
		return cmp.Compare(a.Order, b.Order)
	})
	for _, photo_upload := range photo_uploads {
		fmt.Printf("(%v, %v)\n", photo_upload.Order, photo_upload.File_path)
	}

	return nil
}

func (g *Ghoto) work__upload_photo(
	wg *sync.WaitGroup,
	files_for_worker []File_for_worker,
	ch__photo_upload chan Photo_upload,
) {
	defer wg.Done()

	for _, file_for_worker := range files_for_worker {
		upload_token, err := g.google_photos.Upload_photo(file_for_worker.File_path)
		if err != nil {
			fmt.Printf("❌ Photo upload failed: file=%v\n",
				file_for_worker,
			)
			continue
		}
		ch__photo_upload <- Photo_upload{
			Order:        file_for_worker.Order,
			File_path:    file_for_worker.File_path,
			Upload_token: *upload_token,
		}
	}
}

// func (g *Ghoto) work__create_photo(
// 	worker_id int,
// 	wg *sync.WaitGroup,
// 	files []string,
// 	google_album *googlephotos.Google_album,
// ) {
// 	google_photo, err := g.google_photos.Create_photo(*upload_token, *google_album)
// 	if err != nil {
// 		fmt.Printf("❌ Photo upload failed: #%v-%v, file=%v\n",
// 			worker_id+1,
// 			i+1,
// 			photo_file,
// 		)
// 		continue
// 	}

// 	google_photo__get, get_photo_err := g.google_photos.Get_photo(google_photo.Id)
// 	if get_photo_err == nil && len(google_photo__get.ProductUrl) > 0 {
// 		g.file_remover.Remove(photo_file)

// 		fmt.Printf("✅ Photo upload done: #%v-%v, file=%v, url=%v\n",
// 			worker_id+1,
// 			i+1,
// 			photo_file,
// 			google_photo__get.ProductUrl,
// 		)
// 	} else {
// 		fmt.Printf("❌ Photo upload failed: #%v-%v, file=%v\n",
// 			worker_id+1,
// 			i+1,
// 			photo_file,
// 		)
// 	}
// }

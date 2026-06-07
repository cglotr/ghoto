package ghoto

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func (g *Ghoto) Run(dir string, album string) error {
	dir = filepath.Dir(dir)
	fmt.Printf("🌿 Ghoto %v: dir=%v, album=%v\n",
		ghoto_version,
		dir,
		album,
	)

	files := util.Filter_photo_files(util.Get_files(dir))

	if len(files) < 1 {
		fmt.Printf("✅ No file to process, exiting.\n")
		return nil
	}

	worker_count := min(10, len(files))
	files_per_worker := (len(files) / worker_count) + 1
	if files_per_worker <= 0 {
		files_per_worker = len(files)
	}

	work_assigned_count := 0
	wg := &sync.WaitGroup{}

	for worker_id := range worker_count {
		i := files_per_worker * worker_id
		if i >= len(files) {
			continue
		}
		j := min(i+files_per_worker, len(files))
		files_for_worker := files[i:j]

		wg.Add(1)
		go g.work(
			worker_id,
			wg,
			files_for_worker,
			album,
		)

		work_assigned_count += j - i
		if work_assigned_count >= len(files) {
			break
		}
	}

	wg.Wait()

	files = util.Get_files(dir)
	if len(files) > 0 {
		return errors.New("Run__Get_files__files_remaining")
	}

	non_photo_files := util.Filter_non_photo_files(files)
	for _, non_photo_file := range non_photo_files {
		os.Remove(non_photo_file)
	}

	return nil
}

func (g *Ghoto) work(
	worker_id int,
	wg *sync.WaitGroup,
	files []string,
	album_name string,
) {
	defer wg.Done()

	var google_album *googlephotos.Google_album

	res__list_album, err := g.google_photos.List_album()
	if err != nil {
		fmt.Printf("work__List_album: %v\n", err.Error())
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
			fmt.Printf("work__Create_album: %v\n", err.Error())
		}
	}

	for i, photo_file := range util.Filter_photo_files(files) {
		google_photo, err := g.google_photos.Upload_photo(photo_file, *google_album)
		if err != nil {
			fmt.Printf("work__Upload_photo: %v\n", err.Error())
		}

		google_photo__get, get_photo_err := g.google_photos.Get_photo(google_photo.Id)
		if get_photo_err == nil && len(google_photo__get.ProductUrl) > 0 {
			g.file_remover.Remove(photo_file)

			fmt.Printf("✅ Photo upload done: #%v-%v, file=%v, url=%v\n",
				worker_id+1,
				i+1,
				photo_file,
				google_photo__get.ProductUrl,
			)
		} else {
			fmt.Printf("❌ Photo upload failed: #%v-%v, file=%v\n",
				worker_id+1,
				i+1,
				photo_file,
			)
		}
	}
}

package ghoto

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/cglotr/ghoto/fileremover"
	"github.com/cglotr/ghoto/googlephotos"
	"github.com/stretchr/testify/assert"
)

// --- helpers ---

func makeTempDir(t *testing.T, files []string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range files {
		f, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}
	return dir
}

// --- controllable dummy for Google_photos ---

type fakePhotos struct {
	listAlbums     []googlephotos.Res__album
	listAlbumErr   error
	createAlbumErr error
	uploadPhotoErr error
	getPhotoErr    error
	getPhotoEmpty  bool // returns photo with no ProductUrl
}

func (f *fakePhotos) List_album() (*googlephotos.Res__list_album, error) {
	if f.listAlbumErr != nil {
		return nil, f.listAlbumErr
	}
	return &googlephotos.Res__list_album{Albums: f.listAlbums}, nil
}

func (f *fakePhotos) Create_album(name string) (*googlephotos.Google_album, error) {
	if f.createAlbumErr != nil {
		return nil, f.createAlbumErr
	}
	return &googlephotos.Google_album{Id: "new-album-id", Name: name}, nil
}

func (f *fakePhotos) Upload_photo(file_path string, album googlephotos.Google_album) (*googlephotos.Google_photo, error) {
	if f.uploadPhotoErr != nil {
		return nil, f.uploadPhotoErr
	}
	return &googlephotos.Google_photo{Id: "photo-id", Filename: file_path}, nil
}

func (f *fakePhotos) Get_photo(id string) (*googlephotos.Google_photo, error) {
	if f.getPhotoErr != nil {
		return nil, f.getPhotoErr
	}
	if f.getPhotoEmpty {
		return &googlephotos.Google_photo{Id: id}, nil
	}
	return &googlephotos.Google_photo{Id: id, ProductUrl: "https://photos.google.com/" + id}, nil
}

// --- controllable file remover ---

type fakeRemover struct {
	removed []string
}

func (f *fakeRemover) Remove(file string) error {
	f.removed = append(f.removed, file)
	return os.Remove(file)
}

// --- ghoto factory for tests ---

func newTestGhoto(photos googlephotos.Google_photos, remover fileremover.File_remover) *Ghoto {
	g := Ghoto__new()
	g.google_photos = photos
	g.file_remover = remover
	return g
}

// --- tests ---

func Test__Run(t *testing.T) {
	Ghoto__new().Run("../testfile/", "Album")
}

// Album already exists — photos upload and are removed.
func Test__Run__album_exists(t *testing.T) {
	dir := makeTempDir(t, []string{"a.jpg", "b.mp4", "c.dng", "d.lrv"})
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums: []googlephotos.Res__album{
			{Id: "existing-id", Title: "MyAlbum"},
		},
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "MyAlbum")

	assert.NoError(t, err)
	// jpg and mp4 were uploaded and removed
	assert.Equal(t, 2, len(remover.removed))
	// non-photo files (dng, lrv) cleaned up — they should no longer exist
	assert.NoFileExists(t, filepath.Join(dir, "c.dng"))
	assert.NoFileExists(t, filepath.Join(dir, "d.lrv"))
}

// Album does not exist — Create_album is called.
func Test__Run__album_not_found_creates_album(t *testing.T) {
	dir := makeTempDir(t, []string{"photo.jpg"})
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums: []googlephotos.Res__album{}, // no matching album
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "NewAlbum")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(remover.removed))
}

// No files in directory — Run succeeds with nothing to do.
func Test__Run__empty_dir(t *testing.T) {
	dir := makeTempDir(t, []string{})
	photos := &fakePhotos{
		listAlbums: []googlephotos.Res__album{{Id: "id", Title: "Album"}},
	}
	g := newTestGhoto(photos, &fakeRemover{})

	err := g.Run(dir+"/", "Album")

	assert.NoError(t, err)
}

// Upload fails — file is NOT removed, Run returns error.
func Test__Run__upload_fails(t *testing.T) {
	dir := makeTempDir(t, []string{"photo.jpg"})
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums:     []googlephotos.Res__album{{Id: "id", Title: "Album"}},
		uploadPhotoErr: errors.New("upload failed"),
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "Album")

	assert.Error(t, err)
	assert.Equal(t, 0, len(remover.removed))
}

// Upload succeeds but Get_photo fails — file is NOT removed, Run returns error.
func Test__Run__get_photo_fails(t *testing.T) {
	dir := makeTempDir(t, []string{"photo.jpg"})
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums:  []googlephotos.Res__album{{Id: "id", Title: "Album"}},
		getPhotoErr: errors.New("get photo failed"),
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "Album")

	assert.Error(t, err)
	assert.Equal(t, 0, len(remover.removed))
}

// Get_photo returns a photo without ProductUrl — file is NOT removed.
func Test__Run__get_photo_empty_url(t *testing.T) {
	dir := makeTempDir(t, []string{"photo.jpg"})
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums:    []googlephotos.Res__album{{Id: "id", Title: "Album"}},
		getPhotoEmpty: true,
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "Album")

	assert.Error(t, err)
	assert.Equal(t, 0, len(remover.removed))
}

// Multiple files spread across workers.
func Test__Run__multiple_files(t *testing.T) {
	names := []string{
		"1.jpg", "2.jpg", "3.jpg", "4.jpg", "5.jpg",
		"6.jpg", "7.jpg", "8.jpg", "9.jpg", "10.jpg", "11.jpg",
	}
	dir := makeTempDir(t, names)
	remover := &fakeRemover{}
	photos := &fakePhotos{
		listAlbums: []googlephotos.Res__album{{Id: "id", Title: "Album"}},
	}
	g := newTestGhoto(photos, remover)

	err := g.Run(dir+"/", "Album")

	assert.NoError(t, err)
	assert.Equal(t, len(names), len(remover.removed))
}

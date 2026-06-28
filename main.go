package main

import (
	"flag"
	"fmt"

	p__ghoto "github.com/cglotr/ghoto/ghoto"
	"github.com/cglotr/ghoto/util"
)

func main() {
	f__album := flag.String("album", "Insta360", "Album name to upload to.")
	f__dir := flag.String("dir", "", "Photo directory to upload.")
	f__dryrun := flag.Bool("dryrun", false, "Non-mutative dry run mode.")
	flag.Parse()

	ghoto := p__ghoto.Ghoto__new()
	if !*f__dryrun {
		ghoto.Activate()
	}

	err := ghoto.Run(*f__dir, *f__album)
	if err != nil {
		fmt.Printf("❌ Failed to finish the run: dir=%v, album=%v\n", *f__dir, *f__album)
	}

	files := util.Filter_photo_files(util.Get_files(*f__dir))
	if len(files) > 0 {
		fmt.Printf("❌ Photo files failed to upload: files=%v\n", len(files))
	} else {
		fmt.Printf("🚀 All photo files uploaded!\n")
	}
}

package oakspa

import (
	"io/fs"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(fs fs.FS) oakhttp.Handler {

	return nil
}

// //go:embed all:assets
// var frontend embed.FS
//
// type spaFS struct {
// 	http.FileSystem
// }
//
// func (s *spaFS) Open(name string) (f http.File, err error) {
// 	f, err = s.FileSystem.Open(name)
// 	if errors.Is(err, os.ErrNotExist) && path.Ext(name) == "" {
// 		// file does not exist and the extension is not specified
// 		// the SPA should serve
// 		return s.FileSystem.Open("/index.html")
// 	}
// 	return
// }
//
// func New() http.Handler {
// 	sub, err := fs.Sub(frontend, "assets")
// 	if err != nil {
// 		panic(fmt.Errorf("file system error: %w", err))
// 	}
// 	return http.StripPrefix("/beta/", http.FileServer(&spaFS{FileSystem: http.FS(sub)}))
// }

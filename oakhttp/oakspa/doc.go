/*

Package oakspa adapts an [fs.FS] file system containing frontent application files as an [oakhttp.Handler]. When a request path returns [os.ErrNotExist], a default routing page is served instead. This allows the frontend application to mount browser navigatio and present the correct view.

*/
package oakspa

package static

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

func GetRoutes(root http.FileSystem, path string) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/*", handleFileServer(r, root, path))
	}
}

func handleFileServer(r chi.Router, root http.FileSystem, path string) http.HandlerFunc {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	path += "*"

	return func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}
}

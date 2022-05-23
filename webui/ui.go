package webui

import (
	"embed"
	"errors"
	"fmt"
	"image/jpeg"
	"io/fs"
	"net/http"

	"github.com/areknoster/attendgo/domain"
	"github.com/go-chi/chi/v5"
)

var _ Handler = (*UIHandler)(nil)

type UIHandler struct {
	assets       fs.FS
	photoStorage domain.PhotoStorage
}

func NewUIHandler(photoStorage domain.PhotoStorage) (*UIHandler, error) {
	fs, err := fs.Sub(uiAssets, "ui/public")
	if err != nil {
		return nil, fmt.Errorf("open ui assets: %w", err)
	}
	return &UIHandler{
		assets:       fs,
		photoStorage: photoStorage,
	}, nil
}

//go:embed ui/public/**
var uiAssets embed.FS

// Register implements Handler
func (u *UIHandler) Register(r chi.Router) {
	r.Get("/*", http.FileServer(http.FS(u.assets)).ServeHTTP)
	r.Get("/photos/{photoRef}.jpg", func(w http.ResponseWriter, r *http.Request) {
		rawPhotoRef := chi.URLParam(r, "photoRef")
		photoRef, err := domain.ParsePhotoRef(rawPhotoRef)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		photo, err := u.photoStorage.Get(photoRef)
		if errors.Is(err, domain.ErrDoesntExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "image/jpeg")
		err = jpeg.Encode(w, photo.Img, &jpeg.Options{
			Quality: 80,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

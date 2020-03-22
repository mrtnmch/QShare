package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mxmxcz/qshare/pkg/persistence"
	"github.com/mxmxcz/qshare/pkg/qr"
	"math/rand"
	"net/http"
	"time"
)

const (
	Timeout           = 10 * time.Second
	PathGenerateQR    = "/qr"
	PathUploadContent = "/upload"
	PathFetchContent  = "/fetch"
)

type UploadRequest struct {
	Code     string    `json:"code"`
	Content  string    `json:"content"`
	Uploaded time.Time `json:"-"`
}

type QRResponse struct {
	Code        string `json:"code"`
	Base64Image string `json:"base64image"`
}

type FetchRequest struct {
	Code string `json:"code"`
}

func GetRoutes(generator qr.Generator, store persistence.Manager) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get(PathGenerateQR, handleQRGenerator(generator))
		r.Post(PathUploadContent, handleUpload(generator, store))
		r.Post(PathFetchContent, handleFetch(store))
	}
}

func handleFetch(store persistence.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fetch := FetchRequest{}
		err := render.DecodeJSON(r.Body, &fetch)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, ok := store.Get(fetch.Code)

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		marshal, _ := json.Marshal(result)
		_, err = w.Write(marshal)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		store.Remove(fetch.Code)
	}
}

func handleUpload(generator qr.Generator, store persistence.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upload := UploadRequest{}

		err := render.DecodeJSON(r.Body, &upload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err2 := generator.Validate(upload.Code)

		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		store.Add(persistence.Item{
			Code:    upload.Code,
			Content: upload.Content,
			Created: time.Time{},
			Timeout: Timeout,
		})

		w.WriteHeader(http.StatusAccepted)
	}
}

func handleQRGenerator(generator qr.Generator) http.HandlerFunc {
	rand.Seed(time.Now().UnixNano())

	return func(w http.ResponseWriter, r *http.Request) {
		img, err := generator.Generate()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := &QRResponse{
			Code:        img.Code,
			Base64Image: img.Base64Image,
		}

		marshal, _ := json.Marshal(&response)
		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write(marshal)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

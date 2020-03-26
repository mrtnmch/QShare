package api

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/mxmxcz/qshare/pkg/qr"
	"github.com/mxmxcz/qshare/pkg/repository"
	"github.com/mxmxcz/qshare/pkg/secret"
	"math/rand"
	"net/http"
	"time"
)

const (
	Timeout           = 10 * time.Second
	PathGenerateQR    = "/qr"
	PathUploadContent = "/upload"
	PathFetchContent  = "/fetch"
	HeaderKey         = "X-QShare-Key"
)

type UploadRequest struct {
	Key      string    `json:"code"`
	Content  string    `json:"content"`
	Uploaded time.Time `json:"-"`
}

type QRResponse struct {
	Code        string `json:"code"`
	Base64Image string `json:"base64image"`
}

type FetchRequest struct {
	Key string `json:"code"`
}

type FetchResponse struct {
	Content string `json:"content"`
}

func GetRoutes(imageGenerator qr.Generator, secretProvider secret.Provider, validator secret.Validator, envelopeRepository repository.EnvelopeRepository) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get(PathGenerateQR, handleQRGenerator(imageGenerator, secretProvider))
		r.Post(PathUploadContent, handleUpload(validator, envelopeRepository))
		r.Post(PathFetchContent, handleFetch(envelopeRepository))
	}
}

func handleFetch(store repository.EnvelopeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fetch := FetchRequest{}
		err := render.DecodeJSON(r.Body, &fetch)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result := store.Get(fetch.Key)

		if result == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		marshal, _ := json.Marshal(FetchResponse{string(result.Content)})
		_, err = w.Write(marshal)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		store.Remove(fetch.Key)
	}
}

func handleUpload(validator secret.Validator, store repository.EnvelopeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upload := UploadRequest{}

		err := render.DecodeJSON(r.Body, &upload)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !validator.IsValid([]byte(upload.Key)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		store.Add(&repository.Envelope{
			Key:     upload.Key,
			Content: []byte(upload.Content),
			Created: time.Time{},
			Timeout: Timeout,
		})

		w.WriteHeader(http.StatusAccepted)
	}
}

func handleQRGenerator(imageGenerator qr.Generator, provider secret.Provider) http.HandlerFunc {
	rand.Seed(time.Now().UnixNano())

	return func(w http.ResponseWriter, r *http.Request) {
		key, err := provider.Generate()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		png, err := imageGenerator.NewPNGImage(key)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "image/png")
		w.Header().Add(HeaderKey, string(key))
		_, err = w.Write(png.Content)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

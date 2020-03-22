package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mxmxcz/qshare/pkg"
	"github.com/mxmxcz/qshare/pkg/api"
	"github.com/mxmxcz/qshare/pkg/persistence"
	"github.com/mxmxcz/qshare/pkg/qr"
	"github.com/mxmxcz/qshare/pkg/static"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	ServerTimeout = 60 * time.Second
	SecretLength  = 100
	RsaKeyLength  = 2048
	ServerPort    = 8080
	StaticDir     = "static"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(ServerTimeout))

	secret := generateSecret(SecretLength)
	privateKey, err := rsa.GenerateKey(rand.Reader, RsaKeyLength)
	generator := qr.NewGenerator(secret, privateKey)
	store := persistence.NewMemoryItemManager()

	if err != nil {
		log.Fatal(err)
	}

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, StaticDir)

	r.Route("/", static.GetRoutes(http.Dir(filesDir), "/"))
	r.Route("/api", api.GetRoutes(generator, store))

	log.Printf("Listening on port %d...\n", ServerPort)

	if err = http.ListenAndServe(":"+strconv.Itoa(ServerPort), r); err != nil {
		log.Fatal(err)
	}

	store.Close()
}

func generateSecret(len int) string {
	return pkg.GenerateRandomString(len)
}

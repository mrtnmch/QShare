package main

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mxmxcz/qshare/pkg/api"
	"github.com/mxmxcz/qshare/pkg/qr"
	"github.com/mxmxcz/qshare/pkg/random"
	"github.com/mxmxcz/qshare/pkg/repository"
	"github.com/mxmxcz/qshare/pkg/secret"
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

	globalSecret := generateSecret(SecretLength)
	privateKey, err := rsa.GenerateKey(rand.Reader, RsaKeyLength)

	if err != nil {
		log.Fatal(err)
	}

	encryptor := secret.NewRSAEncryptor(&privateKey.PublicKey, globalSecret)
	decryptor := secret.NewRSADecryptor(privateKey)
	validator := secret.NewBase64SecretValidator(globalSecret, decryptor)
	secretProvider := secret.NewBase64Provider(secret.NewStaticProvider(encryptor, globalSecret))

	imageGenerator := qr.NewGenerator()
	envelopeRepository := repository.NewInMemoryEnvelopeRepository()

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, StaticDir)

	r.Route("/", static.GetRoutes(http.Dir(filesDir), "/"))
	r.Route("/api", api.GetRoutes(imageGenerator, secretProvider, validator, envelopeRepository))

	log.Printf("Listening on port %d...\n", ServerPort)

	if err = http.ListenAndServe(":"+strconv.Itoa(ServerPort), r); err != nil {
		log.Fatal(err)
	}
}

func generateSecret(len int) []byte {
	return random.NewGenerator()(len)
}

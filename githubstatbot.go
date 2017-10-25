package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"golang.org/x/crypto/acme/autocert"
)

//For run:
//env PORT=8080 DB_PATH=/data/githubstatbot/boltdb.db GITHUB_CLIENT_ID= GITHUB_CLIENT_SECRET= TELEGRAM_TOKEN= go run githubstatbot.go
func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		log.Panic("Port is empty")
	}

	tlsDir := os.Getenv("TLS_DIR")

	logFile := os.Getenv("LOG_FILE")
	configureLog(logFile)

	path := os.Getenv("DB_PATH")
	if path == "" {
		log.Panic("DB path is empty")
	}

	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Panic("ClientId or clientSecret is empty")
	}

	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		log.Panic("Telegram token is empty")
	}

	username := os.Getenv("BA_USERNAME")
	password := os.Getenv("BA_PASSWORD")
	if username == "" || password == "" {
		log.Panic("Credential for basic auth is incorrect")
	}

	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		log.Panic("Path for static content is empty")
	}

	db := storage.New(path)
	stateStore := storage.NewStateStore()
	oAuth := github.NewOAuth(clientID, clientSecret)

	bot, err := telegram.NewBot(telegramToken, false, db, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}
	go bot.ReadUpdates()

	basicAuth := &api.BasicAuth{Username: username, Password: password}

	handler := api.New(oAuth, db, stateStore, bot, basicAuth, staticPath)

	router := httprouter.New()
	router.GET("/", handler.Index)
	router.GET("/version", handler.Version)
	router.GET("/github_redirect", handler.GitHubRedirect)

	//Run HTTPS server
	startHttpsServer(router, tlsDir)
	//Run HTTP server
	fmt.Printf("Starting HTTP server on port %s\n", port)
	http.ListenAndServe(":"+port, http.HandlerFunc(handler.RedirectToHttps))
	//http.ListenAndServe(":"+port, router)
}

func configureLog(logFileAddr string) {
	var path string
	if logFileAddr != "" {
		path = logFileAddr
	} else {
		path = "log.txt"
	}

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func startHttpsServer(h http.Handler, tlsDir string) {
	if tlsDir == "" {
		log.Printf("TLS_DIR is empty, so skip serving https")
		return
	}

	httpsServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      h,
	}

	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache(tlsDir),
	}

	httpsServer.Addr = ":443"
	httpsServer.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	go func() {
		fmt.Printf("Starting HTTPS server on %s\n", httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		}
	}()
}

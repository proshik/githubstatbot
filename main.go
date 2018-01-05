package main

import (
	"crypto/tls"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/config"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	//configureLog(cfg.LogDir)

	db := storage.New(cfg.DbPath)
	stateStore := storage.NewStateStore()
	oAuth := github.NewOAuth(cfg.GitHubClientId, cfg.GitHubClientSecret)

	bot, err := telegram.NewBot(cfg.TelegramToken, false, db, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}
	go bot.ReadUpdates()

	basicAuth := &api.BasicAuth{Username: cfg.AuthBasicUsername, Password: cfg.AuthBasicPassword}

	handler := api.New(oAuth, db, stateStore, bot, basicAuth, cfg.StaticFilesDir)

	router := httprouter.New()
	router.GET("/", handler.Index)
	router.GET("/version", handler.Version)
	router.GET("/github_redirect", handler.GitHubRedirect)

	//Run HTTPS server
	log.Printf("Starting HTTP server in mode %s on port %s\n", cfg.Mode, cfg.Port)
	if cfg.Mode == "local" {
		http.ListenAndServe(":"+cfg.Port, router)
	} else {
		startHttpsServer(router, cfg.TlsDir)
		http.ListenAndServe(":"+cfg.Port, http.HandlerFunc(handler.RedirectToHttps))
	}
}

func configureLog(logFileAddr string) {
	if logFileAddr == "" {
		panic(errors.New("Log file is empty"))
	}

	logFilePath := logFileAddr + "/githubstatbot.log"

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
		file.Chmod(0755)
		file.Close()
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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
		log.Printf("Starting HTTPS server on %s\n", httpsServer.Addr)
		err := httpsServer.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		}
	}()
}

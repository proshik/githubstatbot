package main

import (
	"crypto/tls"
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/api"
	"github.com/proshik/githubstatbot/config"
	"github.com/proshik/githubstatbot/github"
	"github.com/proshik/githubstatbot/storage"
	"github.com/proshik/githubstatbot/telegram"
	"golang.org/x/crypto/acme/autocert"
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
	// config logging
	log.SetOutput(os.Stdout)
	// init connect to db(boltDB)
	db := storage.New(cfg.DbPath)
	// create storage for generated statuses for request to github.com
	stateStore := storage.NewStateStore()
	// create oAuth object
	oAuth := github.NewOAuth(cfg.GitHubClientId, cfg.GitHubClientSecret)
	// create Telegram Bot object
	bot, err := telegram.NewBot(cfg.TelegramToken, false, db, stateStore, oAuth)
	if err != nil {
		log.Panic(err)
	}
	// run major method for read updates messages from telegram
	go bot.ReadUpdates()

	// initialize handler
	basicAuth := &api.BasicAuth{Username: cfg.AuthBasicUsername, Password: cfg.AuthBasicPassword}
	handler := api.New(oAuth, db, stateStore, bot, basicAuth, cfg.StaticFilesDir)
	// configuration router
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
	m.HTTPHandler(h)

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

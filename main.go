package main

import (
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
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
	"path"
	"time"
)

type Specification struct {
	Mode     string
	Name     string
	URI      string
	Profile  string
	Username string
	Password string
}

func main() {
	var s Specification
	err := envconfig.Process("githubstatbot", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	springConfig := config.NewSpringConfig(s.Name, s.URI, s.Profile, s.Username, s.Password)
	config, err := springConfig.Read()
	if err != nil {
		panic(err)
	}

	port := config["port"]
	tlsDir := config["tls-dir"]
	configureLog(path.Join(config["log-dir"], "githubstatbot.log"))
	dbPath := config["db-path"]
	clientID := config["github.client-id"]
	clientSecret := config["github.client-secret"]
	telegramToken := config["telegram.token"]
	username := config["authentificate.basic.username"]
	password := config["authentificate.basic.password"]
	staticPath := config["static-dir"]

	db := storage.New(dbPath)
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
	if s.Mode == "local" {
		http.ListenAndServe(":"+port, router)
	} else {
		startHttpsServer(router, tlsDir)
		//Run HTTP server
		fmt.Printf("Starting HTTP server on port %s\n", port)
		http.ListenAndServe(":"+port, http.HandlerFunc(handler.RedirectToHttps))
	}
}

func configureLog(logFileAddr string) {
	if logFileAddr == "" {
		panic(errors.New("Log file is empty"))
	}

	logFile, err := os.OpenFile(logFileAddr, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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

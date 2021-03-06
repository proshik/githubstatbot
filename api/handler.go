package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/proshik/githubstatbot/telegram"
	"io"
	"log"
	"net/http"
	"strings"
)

type AccessTokenReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

var client = &http.Client{}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("Request on index.html")
	http.ServeFile(w, r, h.staticPath+"/index.html")
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Printf("Request on /version")
	if checkAuth(w, r, h.basicAuth) {
		io.WriteString(w, "<html><body>Version: 0.5.3</body></html>")
		return
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func checkAuth(_ http.ResponseWriter, r *http.Request, ba *BasicAuth) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == ba.Username && pair[1] == ba.Password
}

func (h *Handler) GitHubRedirect(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Printf("Unexpected behavior. Code is empty in response from github")
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Printf("Unexpected behavior. State is empty in response from github")
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}

	//check state on valid and get chatId by state
	chatId, err := h.stateStore.Get(state)
	if err != nil {
		log.Printf("Not found chatId by state. Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//delete chatId value state store
	h.stateStore.Delete(state)

	//Build request for get accessToken
	bodyReq := AccessTokenReq{h.oAuth.ClientId, h.oAuth.ClientSecret, code}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bodyReq)

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", b)
	if err != nil {
		log.Printf("Erorr on build request object. Error: %v\n", err)
		h.bot.InformAuth(chatId, false)
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erorr on request got github.com/login/oauth/access_token. %v\n", err)
		h.bot.InformAuth(chatId, false)
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}

	defer resp.Body.Close()

	//decode response with accessToken
	var bodyResp AccessTokenResp
	json.NewDecoder(resp.Body).Decode(&bodyResp)

	if bodyResp.AccessToken == "" {
		log.Printf("Unexpected error. AccessToken from github is empty")
		h.bot.InformAuth(chatId, false)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//save token in storage
	err = h.tokenStore.Add(int64(chatId), bodyResp.AccessToken)
	if err != nil {
		log.Printf("Error on add GitHub user token in db, %v\n", err)
		h.bot.InformAuth(chatId, false)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//inform user in bot about success auth
	h.bot.InformAuth(chatId, true)
	//redirect user to bot page in telegram
	log.Printf("Was authentication user with chatId=%d", int64(chatId))
	http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
}

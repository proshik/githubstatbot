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

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(w, r, h.staticPath+"/index.html")
}

func (h *Handler) Version(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if checkAuth(w, r, h.basicAuth) {
		io.WriteString(w, "<html><body>Version: 0.4.0</body></html>")
		return
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}

func checkAuth(w http.ResponseWriter, r *http.Request, ba *BasicAuth) bool {
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
		log.Printf("Code is empty in response from github")
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}
	state := r.URL.Query().Get("state")
	if state == "" {
		log.Printf("State is empty in response from github")
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
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erorr on build request object. Error: %v\n", err)
		h.bot.InformAuth(chatId, false)
		http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
		return
	}

	defer resp.Body.Close()
	//decode response with accessToken
	var bodyResp AccessTokenResp
	json.NewDecoder(resp.Body).Decode(&bodyResp)

	//save token in storage
	h.tokenStore.Add(int64(chatId), bodyResp.AccessToken)
	//inform user in bot about success auth
	h.bot.InformAuth(chatId, true)
	//redirect user to bot page in telegram
	http.Redirect(w, r, telegram.RedirectBotAddress, http.StatusMovedPermanently)
}

package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
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

func (h *Handler) GitHubAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		log.Printf("Error on received response with code from GitHub.com. Code is empty.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var bodyResp AccessTokenResp
	json.NewDecoder(resp.Body).Decode(&bodyResp)

	fmt.Printf("Received access_token=%s\n", bodyResp.AccessToken)

	chatId, err := strconv.Atoi(state)
	if err != nil {
		log.Printf("Error on convert code=%s to chatId\n", code)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.tokenStore.Add(int64(chatId), bodyResp.AccessToken)

	http.Redirect(w, r, "https://t.me/GitHubStatBot", http.StatusMovedPermanently)
}

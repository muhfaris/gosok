package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muhfaris/gosok/internal/pkg/provider"
)

type (
	ResponseGoogle struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Link       string `json:"link"`
		Picture    string `json:"picture"`
		Locale     string `json:"locale"`
	}
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := provider.AuthUrlGoogle()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, accessToken, err := provider.GetUserInfoGoogle(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var responseGoogle ResponseGoogle
	json.Unmarshal([]byte(content), &responseGoogle)

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//set value
	session.Values["gosok-key"] = accessToken
	session.Values["gosok-id"] = responseGoogle.ID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user", http.StatusFound)
	//fmt.Fprintf(w, "Content:%s\n%s", content, accessToken)
}

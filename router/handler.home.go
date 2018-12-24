package router

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionToken := session.Values["gosok-key"]
	if sessionToken != nil {
		http.Redirect(w, r, "/user", http.StatusFound)
	}

	err = rnd.HTML(w, http.StatusOK, "indexPage", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

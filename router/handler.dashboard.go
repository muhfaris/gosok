package router

import (
	"log"
	"net/http"
)

func DashboardUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionToken := session.Values["gosok-key"]
	if sessionToken == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	err = rnd.HTML(w, http.StatusOK, "dashboard", nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Fatal(err)
	}

	//fmt.Fprintf(w, "Dashboard: %s\n", session.Values["gosok-id"])
}

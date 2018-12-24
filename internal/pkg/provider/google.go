package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

// ClaimSet represents an IdToken response.
type ClaimSet struct {
	Sub string
}

var (
	googleConfig      *oauth2.Config
	googleStateString uuid.UUID
)

func AuthUrlGoogle() string {
	gState := googleStateString.String()
	url := googleConfig.AuthCodeURL(gState)
	return url
}

func GetUserInfoGoogle(state string, code string) ([]byte, string, error) {
	if state != googleStateString.String() {
		return nil, "", fmt.Errorf("invalid oauth state")
	}

	token, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, token.AccessToken, nil
}

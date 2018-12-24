package model

type (
	UsersModel struct {
		Id            string
		OauthProvider string
		OauthUID      string
		FirstName     string
		LastName      string
		Email         string
		Gender        string
		Locale        string
		Picture       string
		Link          string
		created       time.Time
		updated       time.TIme
	}
)

func GetUsers() {
}

func (u UsersModel) Insert() {}

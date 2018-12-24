package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/renderer"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/muhfaris/gosok/internal/pkg/provider"
	"github.com/muhfaris/gosok/router"

	_ "github.com/lib/pq"
)

type dbOptions struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
	SslMode  string
}

var (
	cfgfile         string
	dbPool          *sql.DB
	staticDirectory = "./assets/"
)

var rootCmd = &cobra.Command{
	Use:   "gosok",
	Short: "description of application",
	Long:  "long description of application",
	Run: func(cmd *cobra.Command, args []string) {
		initDB()
		initTemplate()
		provider.Init(provider.InitOption{
			GoogleConfig: &oauth2.Config{
				RedirectURL:  viper.GetString("google.redirect_url"),
				Endpoint:     google.Endpoint,
				Scopes:       []string{viper.GetString("google.scopes_login"), viper.GetString("google.scopes_email"), viper.GetString("google.scopes_openid")},
				ClientID:     viper.GetString("google.client_id"),
				ClientSecret: viper.GetString("google.client_secret"),
			},
		})

		r := mux.NewRouter()
		// for static files
		serveStatic(r, staticDirectory)

		r.HandleFunc("/", router.Home)
		r.HandleFunc("/google/login", router.GoogleLogin)
		r.HandleFunc("/google/callback", router.GoogleCallback)

		r.HandleFunc("/user", router.DashboardUser)
		r.HandleFunc("/logout", router.Logout)

		// Listen Application
		log.Println("Application running in: 1212")
		log.Fatal(http.ListenAndServe(":1212", r))
	},
}

func serveStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"styles":    staticDirectory + "/styles/",
		"scripts":   staticDirectory + "/scripts/",
		"images":    staticDirectory + "/images/",
		"bootstrap": staticDirectory + "/bootstrap/",
		"vendor":    staticDirectory + "vendor",
	}

	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, http.FileServer(http.Dir(pathValue))))
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigType("toml")
	if cfgfile != "" {
		viper.SetConfigFile(cfgfile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//Search config with name ".name"
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".config")
	}

	//Read env
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initDB() {
	DBOptions := dbOptions{
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		DBName:   viper.GetString("database.name"),
		SslMode:  viper.GetString("database.sslmode"),
	}

	db, err := connect(DBOptions)
	if err != nil {
		log.Panic("error/database/connect: ", err)
	}

	dbPool = db
}

func initTemplate() {
	opts := renderer.Options{
		ParseGlobPattern: "./web/template/*.html",
	}

	router.Init(
		dbPool,
		router.InitOption{
			Rnd: renderer.New(opts),
			Session: &router.SessionOptions{
				Domain:   viper.GetString("session.domain"),
				Path:     viper.GetString("session.path"),
				MaxAge:   viper.GetInt("session.max_age"),
				HttpOnly: viper.GetBool("session.http_only"),
			},
		},
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func connect(options dbOptions) (*sql.DB, error) {
	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.DBName,
		options.SslMode,
	)

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

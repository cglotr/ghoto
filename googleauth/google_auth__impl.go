package googleauth

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google_auth__impl struct {
	config *oauth2.Config
	client *http.Client
}

func Google_auth__impl__new() *Google_auth__impl {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		panic("Google_auth__impl__new__ReadFile: " + err.Error())
	}
	config, err := google.ConfigFromJSON(
		b,
		"https://www.googleapis.com/auth/photoslibrary.appendonly",
		"https://www.googleapis.com/auth/photoslibrary.readonly.appcreateddata",
	)
	if err != nil {
		panic("Google_auth__impl__new__ConfigFromJSON: " + err.Error())
	}
	config.RedirectURL = "http://localhost:8080/callback"
	return &Google_auth__impl{
		config: config,
	}

}

func (g *Google_auth__impl) Get_client() (*http.Client, error) {
	if g.client == nil {
		g.client = get_web_auth_client(context.Background(), g.config)
	}
	return g.client, nil
}

func get_web_auth_client(ctx context.Context, config *oauth2.Config) *http.Client {
	code_chan := make(chan string)
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseGlob("googleauth/*.html")
		if err != nil {
			panic("get_web_auth_client__ParseGlob: " + err.Error())
		}

		code := r.URL.Query().Get("code")
		w.Header().Set("Content-Type", "text/html")

		err = tmpl.Execute(w, nil)
		if err != nil {
			panic("get_web_auth_client__Execute: " + err.Error())
		}

		code_chan <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Local server failed: %v", err)
		}
	}()

	auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("👉 Log in to authorize via your default browser:\n\n%s\n\n", auth_url)

	auth_code := <-code_chan
	server.Shutdown(ctx)

	tok, err := config.Exchange(ctx, auth_code)
	if err != nil {
		log.Fatalf("Token exchange failed: %v", err)
	}

	return config.Client(ctx, tok)
}

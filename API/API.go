package API

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	config       *oauth2.Config
	randomString = "random"
	mux          map[string]func(http.ResponseWriter, *http.Request)
)

func init() {

	config = &oauth2.Config{
		RedirectURL:  "http://localhost:8585/callback",
		ClientID:     "ClientID must get from google",
		ClientSecret: "ClientSecret must get from google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	mux["/"] = func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "./Page/index.html")
	}
	mux["/signIn"] = func(writer http.ResponseWriter, request *http.Request) {

		http.Redirect(writer, request, config.AuthCodeURL(randomString), http.StatusTemporaryRedirect)
	}
	mux["/callBack"] = func(writer http.ResponseWriter, request *http.Request) {
		c, err := func(state string, code string) ([]byte, error) {
			if state != randomString {
				return nil, fmt.Errorf("invalid state")
			}

			t, err := config.Exchange(oauth2.NoContext, code)
			if err != nil {
				return nil, fmt.Errorf("exchange failed: %s", err.Error())
			}

			r, err := http.Get(fmt.Sprintf("%s%s", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=", t.AccessToken))
			if err != nil {
				return nil, fmt.Errorf("user info failed: %s", err.Error())
			}

			defer r.Body.Close()
			contents, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, fmt.Errorf("failed reading response body: %s", err.Error())
			}

			return contents, nil
		}(request.FormValue("state"), request.FormValue("code"))
		if err != nil {
			http.Redirect(writer, request, "/", http.StatusNotFound)
			return
		}

		fmt.Fprintf(writer, "content: %s\n", c)
	}
}

type RequestHandler struct{}

func (*RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if method, ok := mux[r.URL.Path]; ok {
		method(w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
	http.Error(w, "I dont know what's happened!!!", http.StatusBadRequest)

}

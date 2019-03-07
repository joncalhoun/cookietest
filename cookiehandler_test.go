package cookies

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCookieHandler(t *testing.T) {
	setHandler := func(name, value string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookieHandler := &CookieHandler{}
			cookieHandler.Set(w, name, value)
		}
	}
	getHandler := func(t *testing.T, name, want string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookieHandler := &CookieHandler{}
			got := cookieHandler.Get(r, name)
			if got != want {
				t.Errorf("Cookie = %v; want %v", got, want)
			}
		}
	}

	var testTable = []struct {
		cookieName  string
		cookieValue string
	}{
		{
			cookieName:  "test",
			cookieValue: "testvalue",
		},
		{
			cookieName:  "dog",
			cookieValue: "cat",
		},
	}
	for i, tt := range testTable {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("/set", setHandler(tt.cookieName, tt.cookieValue))
			mux.Handle("/get", getHandler(t, tt.cookieName, tt.cookieValue))
			server := httptest.NewServer(mux)
			defer server.Close()

			jar, err := cookiejar.New(nil)
			if err != nil {
				t.Fatalf("cookiejar.New() err = %v; want %v", err, nil)
			}
			client := http.Client{
				Jar: jar,
			}
			_, err = client.Get(server.URL + "/set")
			if err != nil {
				t.Fatalf("Get() err = %v; want %v", err, nil)
			}
			_, err = client.Get(server.URL + "/get")
			if err != nil {
				t.Fatalf("Get() err = %v; want %v", err, nil)
			}
		})
	}
}

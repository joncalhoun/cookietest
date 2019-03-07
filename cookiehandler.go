package cookies

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
)

var (
	// create secure cookie (http://www.gorillatoolkit.org/pkg/securecookie)_
	hashKey  = []byte(securecookie.GenerateRandomKey(32))
	blockKey = []byte(securecookie.GenerateRandomKey(32))
	sc       = securecookie.New(hashKey, blockKey)
)

type CookieHandler struct{}

func (c CookieHandler) Set(w http.ResponseWriter, name, value string) {
	if encoded, err := sc.Encode(name, value); err == nil {
		c := &http.Cookie{
			Name:  name,
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, c)
	} else {
		log.Println(err.Error())
	}
}

func (c CookieHandler) Get(r *http.Request, name string) string {
	var value string
	if c, err := r.Cookie(name); err == nil {
		_ = sc.Decode(name, c.Value, &value)
	}
	return value
}

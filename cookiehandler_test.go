package cookies

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCookieHandler_THISDOESNOTWORK(t *testing.T) {
	var testTable = []struct {
		cookieName  string
		cookieValue string
	}{
		{
			cookieName:  "test",
			cookieValue: "testvalue",
		},
	}

	for i, tt := range testTable {
		cookieHandler := &CookieHandler{}
		w := httptest.NewRecorder()
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cookieHandler.SetCookieHandler(w, tt.cookieName, tt.cookieValue)
			req := &http.Request{Header: w.Header()}

			got := cookieHandler.ReadCookieHandler(req, tt.cookieName)

			if got != tt.cookieValue {
				t.Errorf("cookieHandler.ReadCookieHandler(req, tt.cookieName) got = %s; want = %s", got, tt.cookieValue)
			}
		})
	}
}

func TestCookieHandler_THISWORKS(t *testing.T) {
	var testTable = []struct {
		cookies []http.Cookie
	}{
		{
			cookies: []http.Cookie{
				{
					Name:  "test",
					Value: "testvalue",
				},
			},
		},
		{
			cookies: []http.Cookie{
				{
					Name:  "test1",
					Value: "testvalue1",
				},
				{
					Name:  "test2",
					Value: "testvalue2",
				},
				{
					Name:  "test3",
					Value: "testvalue3",
				},
			},
		},
	}

	for i, tt := range testTable {
		// DO NOT REMOVE THIS!
		// need to create a local copy of the tt for the closure below because running tests in parallel
		tt := tt

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			cookieHandler := &CookieHandler{}
			w := httptest.NewRecorder()
			// add all the cookies
			for _, c := range tt.cookies {
				cookieHandler.SetCookieHandler(w, c.Name, c.Value)
			}

			// fetch all the cookies and verify the values are expected
			for _, c := range tt.cookies {
				want := c.Value
				// Do cookies move to "Cookie" header key when redirecting?
				req := &http.Request{Header: http.Header{"Cookie": w.Header()["Set-Cookie"]}}
				got := cookieHandler.ReadCookieHandler(req, c.Name)
				if got != want {
					t.Errorf("cookieHandler.ReadCookieHandler(req, tt.cookieName) got = %s; want = %s", got, want)
				}
			}
		})
	}
}

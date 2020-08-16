package td

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

const APIPath = "https://api.tdameritrade.com/v1/"

var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://auth.tdameritrade.com/auth",
	TokenURL: APIPath + "oauth2/token",
}

func authServer(ctx context.Context, state, addr string) (resp chan string, err error) {
	var u *url.URL
	if u, err = url.Parse(addr); err != nil {
		return
	}

	srv := http.Server{
		Addr: u.Host,
	}

	resp = make(chan string, 1)
	srv.Handler = http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		resp <- q.Get("code")
		wr.Write([]byte("OK\n"))
		time.AfterFunc(time.Second, func() { srv.Shutdown(ctx) })
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http listen error: %v", err)
		}
		close(resp)
	}()

	return
}

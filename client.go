package td // import "go.oneofone.dev/td"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

const Version = "v0.1"

var (
	ErrNoSymbols         = errors.New("must provide at least one symbol")
	ErrMissingAccountID  = errors.New("missing account id and DefaultAccountID isn't set")
	ErrMissingAuthParams = errors.New("token is invalid and no auth params were provided")
)

// NewWithAutoAuth will return a client if the token is valid, otherwise will create a server listening on addr and
// print an auth url.
func NewWithAutoAuth(ctx context.Context, consumerID string, addr string, tok *oauth2.Token) (c *Client, err error) {
	return New(ctx, consumerID, tok, &AuthParams{
		RedirectURL: addr,
		GetCode: func(state, authCodeURL string) (code string, err error) {
			resp, err := authServer(ctx, state, addr)
			if err != nil {
				return "", err
			}
			log.Printf("Visit the URL for the auth dialog: %v\nafter the redirect make sure it starts with http://", authCodeURL)

			return <-resp, nil
		},
	})
}

type AuthParams struct {
	RedirectURL string
	GetCode     func(state, authCodeURL string) (code string, err error)
}

func New(ctx context.Context, consumerID string, tok *oauth2.Token, params *AuthParams) (c *Client, err error) {
	conf := &oauth2.Config{
		ClientID: consumerID + "@AMER.OAUTHAP",
		Endpoint: Endpoint,
	}

	if tok != nil {
		c = &Client{t: tok, ocfg: conf, c: conf.Client(ctx, tok)}
		return
	}

	if params == nil || params.RedirectURL == "" || params.GetCode == nil {
		err = ErrMissingAuthParams
		return
	}
	conf.RedirectURL = params.RedirectURL

	var (
		state = "td" + Version
		url   = conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
		code  string
	)

	if code, err = params.GetCode(state, url); err != nil {
		return
	}

	if tok, err = conf.Exchange(ctx, code, oauth2.AccessTypeOffline); err != nil {
		return
	}

	c = &Client{t: tok, ocfg: conf, c: conf.Client(ctx, tok)}
	return
}

type Client struct {
	t    *oauth2.Token
	ocfg *oauth2.Config
	c    *http.Client

	DefaultAccountID string

	OnRawResponse func(method, url string, req, resp []byte)
}

func (c *Client) Token(ctx context.Context) (*oauth2.Token, error) {
	ts := c.ocfg.TokenSource(context.Background(), c.t)
	return ts.Token()
}

func (c *Client) Request(ctx context.Context, method, ep string, in, out interface{}) error {
	var buf bytes.Buffer
	if in != nil {
		json.NewEncoder(&buf).Encode(in)
	}
	req, _ := http.NewRequestWithContext(ctx, method, APIPath+ep, bytes.NewReader(buf.Bytes()))
	if buf.Len() > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		if c.OnRawResponse != nil {
			c.OnRawResponse(method, ep, buf.Bytes(), b)
		}
		if len(b) == 0 {
			b = []byte("<nil>")
		}
		return xerrors.Errorf("%d (%s %s): %s", resp.StatusCode, method, ep, b)
	}

	if c.OnRawResponse == nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	c.OnRawResponse(method, ep, buf.Bytes(), b)
	return json.NewDecoder(bytes.NewReader(b)).Decode(out)
}

func (c *Client) Quotes(ctx context.Context, symbols ...string) (out map[string]*Quote, err error) {
	if len(symbols) == 0 {
		err = ErrNoSymbols
		return
	}
	err = c.Request(ctx, "GET", "marketdata/quotes?symbol="+strings.Join(symbols, ","), nil, &out)
	return
}

func (c *Client) Quote(ctx context.Context, symbol string) (q *Quote, err error) {
	err = c.Request(ctx, "GET", "marketdata/quotes/"+symbol+"/quotes", nil, &q)
	return
}

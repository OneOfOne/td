package td

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"golang.org/x/xerrors"
)

func (c *Client) Streamer(ctx context.Context, qos int) (s *Streamer, err error) {
	type loginReq struct {
		Credential string `json:"credential"`
		Token      string `json:"token"`
		Version    string `json:"version"`
		QOSLevel   int    `json:"qosLevel"`
	}

	var up *UserPrincipal
	if up, err = c.UserPrincipals(ctx, AllUserPrincipalFields); err != nil {
		return
	}

	acc := up.Accounts[0]
	si := up.StreamerInfo
	tsInMS := si.TokenTimestamp.Time().Unix() * 1000
	creds := url.Values{
		"userid":      {acc.AccountID},
		"token":       {si.Token},
		"company":     {acc.Company},
		"segment":     {acc.Segment},
		"cddomain":    {acc.AccountCdDomainID},
		"usergroup":   {si.UserGroup},
		"accesslevel": {si.AccessLevel},
		"authorized":  {"Y"},
		"timestamp":   {strconv.FormatInt(tsInMS, 10)},
		"appid":       {si.AppID},
		"acl":         {si.Acl},
	}

	var (
		conn *websocket.Conn
		resp *http.Response
	)
	if conn, resp, err = websocket.DefaultDialer.DialContext(ctx, "wss://"+si.StreamerSocketUrl+"/ws", nil); err != nil {
		return
	}

	if resp.StatusCode != 101 {
		err = xerrors.Errorf("%d: %s", resp.StatusCode, resp.Status)
		return
	}

	s = &Streamer{conn: conn, accID: acc.AccountID, appID: si.AppID, key: up.StreamerSubscriptionKeys.Keys[0].Key}

	defer func() {
		if err != nil {
			s.Close()
			s = nil
		}
	}()

	go s.loop()

	err = s.sendRequest(ctx, "ADMIN", "LOGIN", loginReq{
		Credential: creds.Encode(),
		Token:      si.Token,
		Version:    "1.0",
		QOSLevel:   qos,
	})

	return
}

type Streamer struct {
	mux   sync.Mutex
	conn  *websocket.Conn
	accID string
	appID string
	key   string
	reqID int64
	m     sync.Map

	OnData     func(data []Any)
	OnResponse func(code int, message string)
}

func (s *Streamer) SetQoS(ctx context.Context, qos int) error {
	type qosReq struct {
		QOSLevel int `json:"qoslevel"`
	}
	if qos < 0 || qos > 5 {
		return xerrors.Errorf("%d is out of range, the range is 0 to 5", qos)
	}
	return s.sendRequest(ctx, "ADMIN", "QOS", qosReq{qos})
}

type StreamRequestParams struct {
	Keys   string `json:"keys"`
	Fields string `json:"fields"`
}

func (s *Streamer) AccountActivity(ctx context.Context) (<-chan Any, error) {
	const svc = "ACCT_ACTIVITY"
	const fields = "0,1,2,3"

	return s.Subscribe(ctx, svc, StreamRequestParams{Keys: s.key, Fields: fields})
}

type ChartType string

const (
	EquityChart  ChartType = "CHART_EQUITY"
	FuturesChart ChartType = "CHART_FUTURES"
	OptionsChart ChartType = "CHART_OPTIONS"
)

func (s *Streamer) Chart(ctx context.Context, chartType ChartType, symbols ...string) (<-chan Any, error) {
	const allFields = "0,1,2,3,4,5,6,7"
	return s.Subscribe(ctx, string(chartType), StreamRequestParams{
		Keys:   strings.Join(symbols, ","),
		Fields: allFields,
	})
}

func (s *Streamer) Subscribe(ctx context.Context, svc string, params StreamRequestParams) (<-chan Any, error) {
	v, _ := s.m.LoadOrStore(svc, make(chan Any, 256))
	if err := s.sendRequest(ctx, svc, "SUBS", params); err != nil {
		return nil, err
	}
	ch := v.(chan Any)
	return ch, nil
}

// Unsubcribe will close any channels listening for svc and try to run
// the UNSUBS command, which fails because, well reasons...
func (s *Streamer) Unsubcribe(ctx context.Context, svc string) error {
	ch, _ := s.m.LoadAndDelete(svc)
	if ch, ok := ch.(chan Any); ok {
		close(ch)
	}

	// this always returns error 21
	// return s.sendRequest(ctx, svc, "UNSUBS", nil)
	return nil // s.sendRequest(ctx, svc, "UNSUBS", nil)
}

func (s *Streamer) Close() error {
	req, _ := s.makeRequest("ADMIN", "LOGOUT", nil)
	s.conn.WriteJSON(req)
	return s.conn.Close()
}

func (s *Streamer) loop() {
	for {
		var sr streamResponse
		if err := s.conn.ReadJSON(&sr); err != nil {
			break
		}
		// j, _ := json.Marshal(&sr)
		// log.Printf("%s", j)

		for _, r := range sr.Response {
			if v, ok := s.m.Load(r.RequestID); ok {
				if ch, ok := v.(chan *streamDataResponse); ok {
					r := r
					ch <- r
					close(ch)
					s.m.Delete(r.RequestID)
				}
			}
			if s.OnResponse != nil {
				s.OnResponse(r.Content.Code, r.Content.Msg)
			}
		}

		for _, d := range sr.Data {
			if d.Service != "" {
				if v, ok := s.m.Load(d.Service); ok {
					if ch, ok := v.(chan Any); ok {
						for _, c := range d.Content {
							ch <- c
						}
					}
				}
			}

			if s.OnData != nil {
				s.OnData(d.Content)
			}
		}
	}
}

func (s *Streamer) makeRequest(service, cmd string, params interface{}) (*streamRequests, string) {
	id := strconv.FormatInt(atomic.AddInt64(&s.reqID, 1), 10)
	return &streamRequests{
		Requests: []streamRequest{
			{
				Service:    service,
				RequestID:  id,
				Command:    cmd,
				Account:    s.accID,
				Source:     s.appID,
				Parameters: params,
			},
		},
	}, id
}

func (s *Streamer) sendRequest(ctx context.Context, service, cmd string, params interface{}) (err error) {
	req, id := s.makeRequest(service, cmd, params)
	ch := make(chan *streamDataResponse, 1)
	s.mux.Lock()
	s.m.Store(id, ch)
	err = s.conn.WriteJSON(req)
	s.mux.Unlock()
	if err != nil {
		return
	}
	select {
	case r := <-ch:
		if c := r.Content; c.Code != 0 {
			err = xerrors.Errorf("error %d: %s", c.Code, c.Msg)
			return
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return
}

type streamRequests struct {
	Requests []streamRequest `json:"requests,omitempty"`
}

type streamRequest struct {
	Service    string      `json:"service"`
	RequestID  string      `json:"requestid"`
	Command    string      `json:"command"`
	Account    string      `json:"account"`
	Source     string      `json:"source"`
	Parameters interface{} `json:"parameters"`
}

type streamResponse struct {
	Response []*streamDataResponse `json:"response,omitempty"`
	Data     []*streamSubResponse  `json:"data,omitempty"`
	Notify   []struct {
		Heartbeat string `json:"heartbeat,omitempty"`
	} `json:"notify,omitempty"`
}

type streamResponseBase struct {
	Service   string `json:"service,omitempty"`
	RequestID string `json:"requestid,omitempty"`
	Command   string `json:"command,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type streamDataResponse struct {
	streamResponseBase
	Content struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"msg,omitempty"`
	} `json:"content,omitempty"`
}

type streamSubResponse struct {
	streamResponseBase
	Content []Any `json:"content,omitempty"`
}

package td

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.oneofone.dev/anyx"
)

var nytz *time.Location

func init() {
	var err error
	if nytz, err = time.LoadLocation("America/New_York"); err != nil {
		log.Printf("couldn't load NYC timezone, setting it to time.Local")
		nytz = time.Local
	}
}

func NewYorkTZ() *time.Location { return nytz }

const DateTimeFormat = `2006-01-02T15:04:05+0000`

// DateTime represents a timestamp in TD's format
type DateTime string

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	*dt = DateTime(b)
	return nil
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(dt), nil
}

func (dt DateTime) Time() (t time.Time) {
	const ms = int64(1e12)
	if dt[0] == '"' {
		t, _ = time.Parse(DateTimeFormat, string(dt[1:len(dt)-1]))
	} else {
		n, _ := strconv.ParseInt(string(dt), 10, 64)
		if n > ms {
			n /= 1000
		}
		t = time.Unix(n, 0)
	}
	return t
}

type Strike string

func (s Strike) Value() float64 {
	f, _ := strconv.ParseFloat(string(s), 64)
	return f
}

type OptionChainDate string

func (o OptionChainDate) Time() (t time.Time) {
	if i := strings.LastIndexByte(string(o), ':'); i != -1 {
		t, _ = time.ParseInLocation("2006-01-02", string(o), nytz)
	}

	return
}

func Bool(v bool) *bool {
	return &v
}

func BoolVal(src *bool, def bool) *bool {
	if src == nil {
		return &def
	}
	return src
}

func marshalAny(v Any) (b []byte) {
	b, _ = json.Marshal(v)
	return
}

func valToURL(v anyx.A) (u url.Values) {
	if v == nil {
		return
	}

	anyx.Value(v).ForEach(func(key anyx.A, value anyx.Any) (exit bool) {
		if u == nil {
			u = url.Values{}
		}
		u.Set(key.(string), value.String(true))
		return
	})
	return
}

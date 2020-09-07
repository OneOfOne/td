package td

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/OneOfOne/otk"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	consumerID string
	accountID  string

	ctx = context.Background()
)

func init() {
	log.SetFlags(log.Lshortfile)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	consumerID = os.Getenv("CONSUMER_ID")
	accountID = os.Getenv("ACCOUNT_ID")
}

func Test(t *testing.T) {
	var tok *oauth2.Token
	otk.ReadJSONFile("./.token.json", &tok)
	c, err := NewWithAutoAuth(ctx, consumerID, "http://localhost:9000/", tok)
	if err != nil {
		t.Fatal(err)
	}
	if tok, err = c.Token(ctx); err != nil {
		t.Fatal(err)
	}
	if err := otk.WriteJSONFile("./.token.json", &tok, true); err != nil {
		t.Fatal(err)
	}
	// m, err := c.Quotes(ctx, "GOOG", "AAPL")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("%#+v", m)
	// as, err := c.Accounts(ctx)
	// checkAndPrint(t, as, err)
	// a, err := c.Account(ctx, accountID)
	// checkAndPrint(t, a, err)
	// o, err := c.Orders(ctx, accountID, 0, "", "", "")
	// checkAndPrint(t, o, err)
	// oc, err := c.OptionChain(ctx, "AMD", &OptionChainParams{
	// 	IncludeQuotes: true,
	// 	StrikeCount:   4,
	// 	ExpMonth:      "AUG",
	// })
	// checkAndPrint(t, oc, err)
	// oc, err = c.OptionChain(ctx, "AMD", "strikeCount", "4")
	// checkAndPrint(t, oc, err)
	s, err := c.Streamer(ctx, 0)
	checkAndPrint(t, nil, err)
	defer s.Close()

	ch, err := s.Subscribe(ctx, "QUOTE", StreamRequestParams{Keys: "QQQ,SPY", Fields: "0,1,2,3,4,5,6,7"})
	checkAndPrint(t, nil, err)
	t.Log(<-ch)
	t.Log(<-ch)
	t.Log(<-ch)
	ch, err = s.Chart(ctx, EquityChart, "QQQ", "SPY")
	checkAndPrint(t, nil, err)
	ts := time.Now()
	ach, _ := s.AccountActivity(ctx)
	for i := 0; i < 3; i++ {
		select {
		case v := <-ach:
			t.Log(time.Since(ts), v)
		case v := <-ch:
			t.Log(time.Since(ts), v.Get("7").Time("U").UTC(), v)
		}
		ts = time.Now()
	}

	// t.Log(s.Unsubcribe(ctx, "ACC_ACTIVITY"))
	// t.Log(<-ch)
	// os, err := c.Orders(ctx, "", 0, "", "", "")
	// checkAndPrint(t, nil, err)
	// for _, o := range os {
	// 	t.Logf("%#+v", o)
	// }
	// cs, err := c.PriceHistory(ctx, "AMD", nil)
	// checkAndPrint(t, cs, err)
}

func checkAndPrint(tb testing.TB, i interface{}, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
	if i == nil {
		return
	}
	j, _ := json.MarshalIndent(i, "", "  ")
	tb.Logf("%s", j)
}

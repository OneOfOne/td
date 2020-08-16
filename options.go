package td

import (
	"context"
	"net/url"
	"strings"

	"github.com/OneOfOne/any"
)

type OptionChain struct {
	Symbol              string                `json:"symbol,omitempty"`
	Status              string                `json:"status,omitempty"`
	Underlying          *Underlying           `json:"underlying,omitempty"`
	Strategy            Strategy              `json:"strategy,omitempty"`
	Interval            float64               `json:"interval,omitempty"`
	IsDelayed           bool                  `json:"isDelayed,omitempty"`
	IsIndex             bool                  `json:"isIndex,omitempty"`
	InterestRate        float64               `json:"interestRate,omitempty"`
	UnderlyingPrice     float64               `json:"underlyingPrice,omitempty"`
	Volatility          float64               `json:"volatility,omitempty"`
	DaysToExpiration    float64               `json:"daysToExpiration,omitempty"`
	NumberOfContracts   int                   `json:"numberOfContracts,omitempty"`
	MonthlyStrategyList []MonthlyStrategyList `json:"monthlyStrategyList,omitempty"`
	CallExpDateMap      CallExpDateMap        `json:"callExpDateMap,omitempty"`
	PutExpDateMap       PutExpDateMap         `json:"putExpDateMap,omitempty"`
}

type Underlying struct {
	Symbol            string  `json:"symbol,omitempty"`
	Description       string  `json:"description,omitempty"`
	Change            float64 `json:"change,omitempty"`
	PercentChange     float64 `json:"percentChange,omitempty"`
	Close             float64 `json:"close,omitempty"`
	QuoteTime         int64   `json:"quoteTime,omitempty"`
	TradeTime         int64   `json:"tradeTime,omitempty"`
	Bid               float64 `json:"bid,omitempty"`
	Ask               float64 `json:"ask,omitempty"`
	Last              float64 `json:"last,omitempty"`
	Mark              float64 `json:"mark,omitempty"`
	MarkChange        float64 `json:"markChange,omitempty"`
	MarkPercentChange float64 `json:"markPercentChange,omitempty"`
	BidSize           int     `json:"bidSize,omitempty"`
	AskSize           int     `json:"askSize,omitempty"`
	HighPrice         float64 `json:"highPrice,omitempty"`
	LowPrice          float64 `json:"lowPrice,omitempty"`
	OpenPrice         float64 `json:"openPrice,omitempty"`
	TotalVolume       int     `json:"totalVolume,omitempty"`
	ExchangeName      string  `json:"exchangeName,omitempty"`
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh,omitempty"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow,omitempty"`
	Delayed           bool    `json:"delayed,omitempty"`
}

type Option struct {
	PutCall                string      `json:"putCall,omitempty"`
	Symbol                 string      `json:"symbol,omitempty"`
	Description            string      `json:"description,omitempty"`
	ExchangeName           string      `json:"exchangeName,omitempty"`
	Bid                    float64     `json:"bid,omitempty"`
	Ask                    float64     `json:"ask,omitempty"`
	Last                   float64     `json:"last,omitempty"`
	Mark                   float64     `json:"mark,omitempty"`
	BidSize                int         `json:"bidSize,omitempty"`
	AskSize                int         `json:"askSize,omitempty"`
	BidAskSize             string      `json:"bidAskSize,omitempty"`
	LastSize               int         `json:"lastSize,omitempty"`
	HighPrice              float64     `json:"highPrice,omitempty"`
	LowPrice               float64     `json:"lowPrice,omitempty"`
	OpenPrice              float64     `json:"openPrice,omitempty"`
	ClosePrice             float64     `json:"closePrice,omitempty"`
	TotalVolume            int         `json:"totalVolume,omitempty"`
	TradeDate              interface{} `json:"tradeDate,omitempty"`
	TradeTimeInLong        int64       `json:"tradeTimeInLong,omitempty"`
	QuoteTimeInLong        int64       `json:"quoteTimeInLong,omitempty"`
	NetChange              float64     `json:"netChange,omitempty"`
	Volatility             float64     `json:"volatility,omitempty"`
	Delta                  float64     `json:"delta,omitempty"`
	Gamma                  float64     `json:"gamma,omitempty"`
	Theta                  float64     `json:"theta,omitempty"`
	Vega                   float64     `json:"vega,omitempty"`
	Rho                    float64     `json:"rho,omitempty"`
	OpenInterest           int         `json:"openInterest,omitempty"`
	TimeValue              float64     `json:"timeValue,omitempty"`
	TheoreticalOptionValue float64     `json:"theoreticalOptionValue,omitempty"`
	TheoreticalVolatility  float64     `json:"theoreticalVolatility,omitempty"`
	OptionDeliverablesList Any         `json:"optionDeliverablesList,omitempty"`
	StrikePrice            float64     `json:"strikePrice,omitempty"`
	ExpirationDate         int64       `json:"expirationDate,omitempty"`
	DaysToExpiration       int         `json:"daysToExpiration,omitempty"`
	ExpirationType         string      `json:"expirationType,omitempty"`
	LastTradingDay         int64       `json:"lastTradingDay,omitempty"`
	Multiplier             float64     `json:"multiplier,omitempty"`
	SettlementType         string      `json:"settlementType,omitempty"`
	DeliverableNote        string      `json:"deliverableNote,omitempty"`
	IsIndexOption          Any         `json:"isIndexOption,omitempty"`
	PercentChange          float64     `json:"percentChange,omitempty"`
	MarkChange             float64     `json:"markChange,omitempty"`
	MarkPercentChange      float64     `json:"markPercentChange,omitempty"`
	NonStandard            bool        `json:"nonStandard,omitempty"`
	InTheMoney             bool        `json:"inTheMoney,omitempty"`
	Mini                   bool        `json:"mini,omitempty"`
}

type CallExpDateMap map[string]map[Strike][]*Option

type PutExpDateMap map[string]map[Strike][]*Option

type MonthlyStrategyList struct {
	Month              string               `json:"month,omitempty"`
	Year               int                  `json:"year,omitempty"`
	Day                int                  `json:"day,omitempty"`
	DaysToExp          int                  `json:"daysToExp,omitempty"`
	SecondaryMonth     string               `json:"secondaryMonth,omitempty"`
	SecondaryYear      int                  `json:"secondaryYear,omitempty"`
	SecondaryDay       int                  `json:"secondaryDay,omitempty"`
	SecondaryDaysToExp int                  `json:"secondaryDaysToExp,omitempty"`
	Type               string               `json:"type,omitempty"`
	SecondaryType      string               `json:"secondaryType,omitempty"`
	Leap               bool                 `json:"leap,omitempty"`
	OptionStrategyList []OptionStrategyList `json:"optionStrategyList,omitempty"`
	SecondaryLeap      bool                 `json:"secondaryLeap,omitempty"`
}

type StrategyLeg struct {
	Symbol      string  `json:"symbol,omitempty"`
	PutCallInd  string  `json:"putCallInd,omitempty"`
	Description string  `json:"description,omitempty"`
	Bid         float64 `json:"bid,omitempty"`
	Ask         float64 `json:"ask,omitempty"`
	Range       string  `json:"range,omitempty"`
	StrikePrice float64 `json:"strikePrice,omitempty"`
	TotalVolume float64 `json:"totalVolume,omitempty"`
}

type OptionStrategyList struct {
	PrimaryLeg     *StrategyLeg `json:"primaryLeg,omitempty"`
	SecondaryLeg   *StrategyLeg `json:"secondaryLeg,omitempty"`
	StrategyStrike string       `json:"strategyStrike,omitempty"`
	StrategyBid    float64      `json:"strategyBid,omitempty"`
	StrategyAsk    float64      `json:"strategyAsk,omitempty"`
}

type ContractType string

const (
	PutContracts  ContractType = "PUT"
	CallContracts ContractType = "CALL"
	AllContracts  ContractType = "ALL"
)

type Strategy string

const (
	StrategySingle     Strategy = "SINGLE"
	StrategyAnalytical Strategy = "ANALYTICAL"
	StrategyCovered    Strategy = "COVERED"
	StrategyVertical   Strategy = "VERTICAL"
	StrategyCalendar   Strategy = "CALENDAR"
	StrategyStrangle   Strategy = "STRANGLE"
	StrategyStraddle   Strategy = "STRADDLE"
	StrategyButterfly  Strategy = "BUTTERFLY"
	StrategyCondor     Strategy = "CONDOR"
	StrategyDiagonal   Strategy = "DIAGONAL"
	StrategyCollar     Strategy = "COLLAR"
	StrategyRoll       Strategy = "ROLL"
)

type StrikeRange string

const (
	InTheMoney         StrikeRange = "ITM"
	NearTheMoney       StrikeRange = "NTM"
	OutOfTheMoney      StrikeRange = "OTM"
	StrikesAboveMarket StrikeRange = "SAK"
	StrikesBelowMarket StrikeRange = "SBK"
	StrikesNearMarket  StrikeRange = "SNK"
	AllStrikes         StrikeRange = "ALL"
)

type OptionChainParams struct {
	// Type of contracts to return in the chain. Can be CALL, PUT, or ALL. Default is ALL.
	ContractType ContractType `json:"contractType,omitempty"`

	// The number of strikes to return above and below the at-the-money price.
	StrikeCount int

	// Include quotes for options in the option chain. Can be TRUE or FALSE. Default is FALSE.
	IncludeQuotes bool

	// Passing a value returns a Strategy Chain. Possible values are SINGLE, ANALYTICAL (allows use of the volatility, underlyingPrice, interestRate, and daysToExpiration params to calculate theoretical values), COVERED, VERTICAL, CALENDAR, STRANGLE, STRADDLE, BUTTERFLY, CONDOR, DIAGONAL, COLLAR, or ROLL. Default is SINGLE.
	Strategy Strategy

	// Strike interval for spread strategy chains (see strategy param).
	Interval int

	// Provide a strike price to return options only at that strike price.
	Strike float64

	// 	Returns options for the given range. Possible values are:
	// ITM: In-the-money
	// NTM: Near-the-money
	// OTM: Out-of-the-money
	// SAK: Strikes Above Market
	// SBK: Strikes Below Market
	// SNK: Strikes Near Market
	// ALL: All Strikes
	// Default is ALL.
	Range StrikeRange

	// Only return expirations after this date. For strategies, expiration refers to the nearest term expiration in the strategy.
	// Valid ISO-8601 formats are: 2006-01-02 and 2006-01-02T15:04:05
	FromDate DateTime

	// Only return expirations before this date. For strategies, expiration refers to the nearest term expiration in the strategy.
	// Valid ISO-8601 formats are: 2006-01-02 and 2006-01-02T15:04:05
	ToDate DateTime

	// Volatility to use in calculations. Applies only to ANALYTICAL strategy chains (see strategy param).
	Volatility float64

	// Underlying price to use in calculations. Applies only to ANALYTICAL strategy chains (see strategy param).
	UnderlyingPrice float64

	// Interest rate to use in calculations. Applies only to ANALYTICAL strategy chains (see strategy param).
	InterestRate float64

	// Days to expiration to use in calculations. Applies only to ANALYTICAL strategy chains (see strategy param).
	DaysToExpiration int

	// Return only options expiring in the specified month. Month is given in the three character format.
	// Example: JAN
	// Default is ALL.
	ExpMonth string

	// Type of contracts to return. Possible values are:
	// S: Standard contracts
	// NS: Non-standard contracts
	// ALL: All contracts
	OptionType string
}

// OptionChain is complicated... check https://developer.tdameritrade.com/option-chains/apis/get/marketdata/chains
func (c *Client) OptionChain(ctx context.Context, symbol string, params *OptionChainParams) (out *OptionChain, err error) {
	args := url.Values{}
	args.Set("symbol", symbol)
	if params != nil {
		pa := any.Value(params)
		pa.ForEach(func(key any.A, value any.Any) (exit bool) {
			k := key.(string)
			k = strings.ToLower(k[:1]) + k[1:]
			args.Set(k, value.String(true))
			return
		})
	}

	err = c.Request(ctx, "GET", "marketdata/chains?"+args.Encode(), nil, &out)
	return
}

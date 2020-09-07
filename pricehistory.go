package td

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"time"
)

type PeriodType string

const (
	DayPeriod        = PeriodType("day")
	MonthPeriod      = PeriodType("month")
	YearPeriod       = PeriodType("year")
	YearToDatePeriod = PeriodType("ytd")
)

type FrequencyType string

const (
	MinuteFrequency  = FrequencyType("minute")
	DailyFrequency   = FrequencyType("daily")
	WeeklyFrequency  = FrequencyType("weekly")
	MonthlyFrequency = FrequencyType("monthly")
)

type HistoryParams struct {
	// The type of period to show. Valid values are day, month, year, or ytd (year to date). Default is day.
	PeriodType PeriodType `json:"periodType,omitempty"`

	// The number of periods to show.
	// Example: For a 2 day / 1 min chart, the values would be:
	// period: 2
	// periodType: day
	// frequency: 1
	// frequencyType: min
	// Valid periods by periodType (defaults marked with an asterisk):
	// day: 1, 2, 3, 4, 5, 10*
	// month: 1*, 2, 3, 6
	// year: 1*, 2, 3, 5, 10, 15, 20
	// ytd: 1*
	Period int `json:"period,omitempty"`

	// The type of frequency with which a new candle is formed.
	// Valid frequencyTypes by periodType (defaults marked with an asterisk):
	// day: minute*
	// month: daily, weekly*
	// year: daily, weekly, monthly*
	// ytd: daily, weekly*
	FrequencyType FrequencyType `json:"frequencyType,omitempty"`

	// 	The number of the frequencyType to be included in each candle.

	// Valid frequencies by frequencyType (defaults marked with an asterisk):

	// minute: 1*, 5, 10, 15, 30
	// daily: 1*
	// weekly: 1*
	// monthly: 1*
	Frequency int `json:"frequency,omitempty"`

	// End date as milliseconds since epoch. If startDate and endDate are provided, period should not be provided. Default is previous trading day.
	StartDate time.Time `json:"startDate,omitempty"`

	// Start date as milliseconds since epoch. If startDate and endDate are provided, period should not be provided.
	EndDate time.Time `json:"endDate,omitempty"`

	// true to return extended hours data, false for regular market hours only. Default is true
	// use Bool to easily return a pointer
	NeedExtendedHoursData *bool `json:"needExtendedHoursData,omitempty"`
}

func (p *HistoryParams) Query() string {
	if p == nil {
		return ""
	}
	u := url.Values{}
	if p.PeriodType != "" {
		u.Set("periodType", string(p.PeriodType))
	}
	if p.Period != 0 {
		u.Set("period", strconv.Itoa(p.Period))
	}
	if p.FrequencyType != "" {
		u.Set("frequencyType", string(p.FrequencyType))
	}

	if p.Frequency != 0 {
		u.Set("frequency", strconv.Itoa(p.Frequency))
	}

	if !p.StartDate.IsZero() {
		u.Set("startDate", strconv.FormatInt(p.StartDate.Unix()*1000, 10))
	}

	if !p.EndDate.IsZero() {
		u.Set("endDate", strconv.FormatInt(p.EndDate.Unix()*1000, 10))
	}

	u.Set("needExtendedHoursData", strconv.FormatBool(p.NeedExtendedHoursData == nil || *p.NeedExtendedHoursData))
	log.Println(u)
	return u.Encode()
}

type Candle struct {
	Open     float64  `json:"open,omitempty"`
	High     float64  `json:"high,omitempty"`
	Low      float64  `json:"low,omitempty"`
	Close    float64  `json:"close,omitempty"`
	Volume   int      `json:"volume,omitempty"`
	Datetime DateTime `json:"datetime,omitempty"`
}

func (c *Client) PriceHistory(ctx context.Context, symbol string, params *HistoryParams) (_ Candles, err error) {
	var out struct {
		Candles Candles `json:"candles,omitempty"`
	}
	err = c.Request(ctx, "GET", "marketdata/"+symbol+"/pricehistory?"+params.Query(), nil, &out)
	return out.Candles, err
}

type Candles []Candle

func (c Candles) GroupByClose() []float64 {
	out := make([]float64, len(c))
	for i := range c {
		out[i] = c[i].Close
	}
	return out
}

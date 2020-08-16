package td

import (
	"context"
	"strings"
)

type accountWrapper struct {
	SecuritiesAccount Account `json:"securitiesAccount,omitempty"`
}

func (c *Client) Account(ctx context.Context, accountID string) (out *Account, err error) {
	if accountID == "" {
		accountID = c.DefaultAccountID
	}
	if accountID == "" {
		err = ErrMissingAccountID
		return
	}
	var aw accountWrapper
	err = c.Request(ctx, "GET", "accounts/"+accountID, nil, &aw)
	return
}

func (c *Client) Accounts(ctx context.Context) (out []*Account, err error) {
	var aws []accountWrapper
	err = c.Request(ctx, "GET", "accounts", nil, &aws)
	out = make([]*Account, 0, len(aws))
	for i := range aws {
		out = append(out, &aws[i].SecuritiesAccount)
	}
	return
}

var AllUserPrincipalFields = []string{"streamerConnectionInfo", "streamerSubscriptionKeys", "preferences", "surrogateIds"}

func (c *Client) UserPrincipals(ctx context.Context, fields []string) (out *UserPrincipal, err error) {
	err = c.Request(ctx, "GET", "userprincipals?fields="+strings.Join(fields, ","), nil, &out)
	return
}

package td

import (
	"context"
	"strconv"
	"strings"
)

func (c *Client) CancelOrder(ctx context.Context, accountID, orderID string) error {
	if accountID == "" {
		accountID = c.DefaultAccountID
	}
	if accountID == "" {
		return ErrMissingAccountID
	}
	return c.Request(ctx, "GET", "accounts/"+accountID+"/orders/"+orderID, nil, nil)
}

func (c *Client) Order(ctx context.Context, accountID, orderID string) (out *Order, err error) {
	if accountID == "" {
		accountID = c.DefaultAccountID
	}

	if accountID == "" {
		err = ErrMissingAccountID
		return
	}

	err = c.Request(ctx, "GET", "accounts/"+accountID+"/orders/"+orderID, nil, &out)
	return
}

// Orders returns either all the orders for the given accountID, if accountID and DefaultAccountID are empty, it returns all orders for all accounts.
func (c *Client) Orders(ctx context.Context, accountID string, maxResults int, fromEnteredTime, toEnteredTime string, status Status) (out []*Order, err error) {
	if accountID == "" {
		accountID = c.DefaultAccountID
	}
	var ep string
	if accountID == "" {
		ep = "orders?"
	} else {
		ep = "accounts/" + accountID + "/orders?"
	}

	var opts []string
	if maxResults > 0 {
		opts = append(opts, "maxResults="+strconv.Itoa(maxResults))
	}

	if fromEnteredTime != "" {
		opts = append(opts, "fromEnteredTime="+fromEnteredTime)
	}

	if toEnteredTime != "" {
		opts = append(opts, "toEnteredTime="+toEnteredTime)
	}

	if status != "" {
		opts = append(opts, "status="+string(status))
	}

	err = c.Request(ctx, "GET", ep+strings.Join(opts, "&"), nil, &out)
	return
}

func (c *Client) CreateOrder() {}

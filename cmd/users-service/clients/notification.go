package clients

import (
	"context"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/tracing"
)

func NewNotification(url string) *Notification {
	return &Notification{
		url: url,
		client: &http.Client{
			Transport: &tracing.HTTPTransport{
				Transport: http.DefaultTransport,
			},
		},
	}
}

type Notification struct {
	client *http.Client
	url    string
}

func (c *Notification) Notify(ctx context.Context, userID string) error {
	logger := sharedContext.Logger(ctx)
	q := url.Values{"user_id": {userID}}
	req, err := http.NewRequest(http.MethodGet, c.url+"/notify?"+q.Encode(), nil)
	if err != nil {
		logger.Error("could not create request", zap.Error(err))
		return err
	}

	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		logger.Error("could not send notification request", zap.Error(err))
		return err
	}
	defer func() { _ = res.Body.Close() }()

	return nil
}

package headlessclient

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"golang.org/x/time/rate"
)

type RLHeadlessClient struct {
	Ratelimiter *rate.Limiter
}

func (c *RLHeadlessClient) RetrieveHtml(url string, waitTime time.Duration) (string, error) {

	ctx := context.Background()
	err := c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return "", err
	}

	page := rod.New().MustConnect().MustPage(url).MustWaitStable()
	defer page.MustClose()
	page.Mouse.MustScroll(0, 400)

	// sleep provide needed for javascript on browser to finish
	time.Sleep(waitTime)

	stringHtml, err := page.HTML()
	return stringHtml, err
}

// maxRequest 20, period 10*time.Second means max 20 request every 10 second
func NewClient(maxRequest int, period time.Duration) RLHeadlessClient {

	rateLimit := rate.NewLimiter(rate.Every(period), maxRequest)
	c := RLHeadlessClient{
		Ratelimiter: rateLimit,
	}
	return c
}

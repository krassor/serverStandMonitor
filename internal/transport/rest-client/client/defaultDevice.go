package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	retryCount    int           = 3
	retryDuration time.Duration = 3
)

type DefaultDevice struct {
	httpClient *http.Client
}

func NewDefaultDevice(httpClient *http.Client) DeviceStatusClient {
	return &DefaultDevice{httpClient: httpClient}
}

func (d *DefaultDevice) GetStatus(ctx context.Context, url string) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request with ctx: %w", err)
	}

	var res *http.Response
	for i := 0; i < retryCount; i++ {
		r, e := d.httpClient.Do(req)
		if e != nil {
			err = e
			time.Sleep(time.Second * retryDuration)
			continue
		}
		res = r
	}
	if err != nil {
		return false, fmt.Errorf("failed to perform http request: %w", err)
	}
	//if res.StatusCode != http.StatusOK {
	//	return false, fmt.Errorf("status code is not 200 OK")
	//}
	if res.ContentLength == 0 {
		return false, fmt.Errorf("response body is null")
	}

	return true, nil
}

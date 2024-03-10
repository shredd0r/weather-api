package client

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type LoggingTransport struct {
	Logger    *logrus.Logger
	Transport http.RoundTripper
}

// RoundTrip викликається для кожного запиту, що проходить через клієнт.
func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	t.requestLogging(req)

	res, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	t.responseLogging(res)

	return res, nil
}

func (t *LoggingTransport) requestLogging(req *http.Request) {
	t.Logger.Infof("Request: %s %s", req.Method, req.Host)
	t.Logger.Debugf("Url: %s", req.URL)
}

func (t *LoggingTransport) responseLogging(res *http.Response) {
	t.Logger.Infof("Response: %s", res.Status)

	if res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		bytesBody, err := io.ReadAll(res.Body)

		if err != nil {
			panic("cant parse response body")
		}

		res.Body = io.NopCloser(bytes.NewReader(bytesBody))
		t.Logger.Debugf("Body: %s", string(bytesBody))
	}
}

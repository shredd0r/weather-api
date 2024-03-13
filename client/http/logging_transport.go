package http

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type LoggingTransport struct {
	log       *logrus.Logger
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
	t.log.Infof("Request: %s %s", req.Method, req.Host)
	t.log.Debugf("Url: %s", req.URL)
}

func (t *LoggingTransport) responseLogging(res *http.Response) {
	t.log.Infof("Response: %s", res.Status)
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bytesBody, err := io.ReadAll(res.Body)

		if err != nil {
			panic("cant parse response body")
		}

		res.Body = io.NopCloser(bytes.NewReader(bytesBody))
		t.log.Debugf("Body: %s", string(bytesBody))
	}
}

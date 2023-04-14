package kiotalog

import (
	"bytes"
	"github.com/redhat-developer/app-services-cli/pkg/core/logging"
	"io"
	"net/http"

	"github.com/microsoft/kiota-http-go"
)

// Used for rhoas verbose mode (-v --verbose)
type LoggingHandler struct {
	logger logging.Logger
}

func NewLoggingHandler(l logging.Logger) *LoggingHandler {
	return &LoggingHandler{
		logger: l,
	}
}

func (c *LoggingHandler) Intercept(pipeline nethttplibrary.Pipeline, middlewareIndex int, req *http.Request) (*http.Response, error) {

	c.logger.Info("", req.Method, req.URL)

	resp, err := pipeline.Next(req, middlewareIndex)
	if err != nil {
		return nil, err
	}

	c.logger.Info("", resp.Status)
	var readBody []byte

	readBody, err = io.ReadAll(resp.Body)
	if err == nil {
		bodyString := string(readBody)
		c.logger.Info("%s", bodyString)
	}

	resp.Body.Close()

	// Once read the stream is close, we need to create a new Reader for subsequent Middlewares
	// See https://stackoverflow.com/questions/43021058/golang-read-request-body-multiple-times
	newBodyReadable := bytes.NewBuffer(readBody)
	resp.Body = io.NopCloser(newBodyReadable)

	return resp, nil

}

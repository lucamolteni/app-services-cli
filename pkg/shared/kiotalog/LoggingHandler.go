package kiotalog

import (
	"github.com/redhat-developer/app-services-cli/pkg/core/logging"
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

	return pipeline.Next(req, middlewareIndex)

}

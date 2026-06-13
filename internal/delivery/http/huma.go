package http

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// Default API Responses for existing format { "success": true, "message": "...", "data": ... }
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type BackwardCompatibleError struct {
	status  int
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (e *BackwardCompatibleError) Error() string {
	return e.Message
}

func (e *BackwardCompatibleError) GetStatus() int {
	return e.status
}

func initErrorTransformer() {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		details := make([]string, 0, len(errs))
		for _, err := range errs {
			details = append(details, err.Error())
		}
		detail := message
		if len(details) > 0 {
			detail = details[0]
		}
		return &BackwardCompatibleError{
			status:  status,
			Success: false,
			Message: detail,
		}
	}
}

func defaultConfig() huma.Config {
	config := huma.DefaultConfig("tsctl API", "1.0.0")
	config.OpenAPIPath = "/api/openapi" // Will serve /api/openapi.json
	config.DocsPath = ""                // Disable default stoplight docs
	config.Servers = []*huma.Server{
		{URL: "/"},
	}
	return config
}

func scalarDocsHTML(openAPIURL string) string {
	return `<!doctype html>
<html>
  <head>
    <title>tsctl API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>body { margin: 0; padding: 0; }</style>
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '` + openAPIURL + `',
        theme: 'default'
      })
    </script>
  </body>
</html>`
}

func SetupHuma(r *gin.Engine) huma.API {
	initErrorTransformer()
	api := humagin.New(r, defaultConfig())

	// scalar API doc router
	r.GET("/api/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, scalarDocsHTML("/api/openapi.json"))
	})

	return api
}

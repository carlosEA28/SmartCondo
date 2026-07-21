package server

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>SmartCondo API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *::before, *::after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/docs/openapi.yaml',
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset
      ],
      layout: "BaseLayout",
    });
  </script>
</body>
</html>`

func (s *Server) docsHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, swaggerUIHTML)
}

func (s *Server) registerDocsRoutes(router *gin.Engine) {
	yamlPath := filepath.Join(s.config.Docs.Path, "openapi.yaml")

	router.GET("/", s.docsHandler)
	router.GET("/docs", s.docsHandler)
	router.StaticFile("/docs/openapi.yaml", yamlPath)
}

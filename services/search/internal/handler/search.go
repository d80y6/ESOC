package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"
)

type SearchHandler struct {
	client *opensearch.Client
	logger *zap.Logger
}

func NewSearchHandler(client *opensearch.Client, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		client: client,
		logger: logger,
	}
}

type SearchRequest struct {
	Query string `json:"query"`
	From  int    `json:"from"`
	Size  int    `json:"size"`
}

func (h *SearchHandler) SearchLogs(c echo.Context) error {
	var req SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	tenantID, ok := c.Get("tenant_id").(string)
	if !ok || tenantID == "" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "missing tenant context"})
	}

	if req.Size == 0 {
		req.Size = 10
	}

	// Simple match_all if no query, else use query_string for flexibility
	// ENFORCE TENANT ISOLATION
	var queryBody string
	if req.Query == "" {
		queryBody = fmt.Sprintf(`{"query": {"bool": {"filter": [{"term": {"tenant_id": "%s"}}]}}}`, tenantID)
	} else {
		// Use a bool query to combine user query and tenant filter
		userQuery := strings.ReplaceAll(req.Query, `"`, `\"`)
		queryBody = fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{"query_string": {"query": "%s"}}
					],
					"filter": [
						{"term": {"tenant_id": "%s"}}
					]
				}
			}
		}`, userQuery, tenantID)
	}

	res, err := h.client.Search(
		h.client.Search.WithContext(context.Background()),
		h.client.Search.WithIndex("normalized-events-*"),
		h.client.Search.WithBody(strings.NewReader(queryBody)),
		h.client.Search.WithFrom(req.From),
		h.client.Search.WithSize(req.Size),
		h.client.Search.WithSort("@timestamp:desc"),
	)
	if err != nil {
		h.logger.Error("failed to search opensearch", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "search failed"})
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to decode results"})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *SearchHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/search", h.SearchLogs)
}

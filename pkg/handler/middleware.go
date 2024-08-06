package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	customerCtx         = "customerId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty Authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid Authorization header")
		return
	}

	customerId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set(customerCtx, customerId)

}

func (h *Handler) getId(c *gin.Context) (int, error) {
	id, ok := c.Get(customerCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "customer id not found")
		return 0, errors.New("customer id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "customer id not found")
		return 0, errors.New("customer id not found")
	}
	return idInt, nil
}

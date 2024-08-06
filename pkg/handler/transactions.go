package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type transactionInput struct {
	Value int `json:"value" binding:"required"`
}

func (h *Handler) depositMoney(c *gin.Context) {
	id, err := h.getId(c)
	if err != nil {
		return
	}

	var input transactionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	balance, err := h.services.Transaction.DepositMoney(id, input.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
	})

}

func (h *Handler) withdrawMoney(c *gin.Context) {
	id, err := h.getId(c)
	if err != nil {
		return
	}

	var input transactionInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	balance, err := h.services.Transaction.WithdrawMoney(id, input.Value)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
	})
}

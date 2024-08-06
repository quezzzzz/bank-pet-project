package handler

import (
	"bank"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type creditInput struct {
	Value    int `json:"value" binding:"required"`
	Variable int `json:"variable" binding:"required"`
}

func (h *Handler) takeCredit(c *gin.Context) {
	id, err := h.getId(c)
	var input creditInput
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadGateway, "invalid input")
		return
	}
	var credit bank.Credit
	balance, credit, err := h.services.TakeCredit(id, input.Value, input.Variable)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		time.Sleep(time.Duration(credit.LoanPeriod) * time.Second)
		_, err := h.services.CloseCredit(credit.Id, credit.CustomerId, credit.CurrentDebt)
		if err != nil {
			logrus.Println(err)
		}
		logrus.Println("1")
	}()
	c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
	})

}

//ВОЗВРАЩАТЬ ВСЮ ИНФУ + ОБЯЗ ПЛАТЕЖ + С БИНД

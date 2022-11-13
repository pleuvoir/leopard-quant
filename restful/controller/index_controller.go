package controller

import (
	"github.com/gin-gonic/gin"
	"leopard-quant/restful/model"
	"net/http"
)

func NewIndexController() *IndexController {
	return &IndexController{}
}

type IndexController struct {
}

func (t *IndexController) Index(ct *gin.Context) {
	rs := model.Success("success")
	ct.JSON(http.StatusOK, rs)
}

func (t *IndexController) Welcome(c *gin.Context) {
	success := model.Success("welcome")
	success.SetData("u find new place")
	c.JSON(http.StatusOK, success)
}

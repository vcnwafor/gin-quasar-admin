package public_api

import (
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/global"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/gqaplugin/example/service/public_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetNews(c *gin.Context)  {
	err, news := public_service.GetNews()
	if err != nil {
		global.GqaLog.Error("获取新闻列表！", zap.Any("err", err))
		global.ErrorMessage("获取新闻列表，"+err.Error(), c)
	} else {
		global.SuccessMessageData(gin.H{"records": news}, "获取新闻列表！", c)
	}
}
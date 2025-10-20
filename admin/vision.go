package admin

import (
	"chat/channel"
	"chat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VisionConfigForm struct {
	Models []string `json:"models"`
}

// GetVisionConfigAPI 获取视觉模型配置
func GetVisionConfigAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   channel.SystemInstance.Vision,
	})
}

// UpdateVisionConfigAPI 更新视觉模型配置
func UpdateVisionConfigAPI(c *gin.Context) {
	var form VisionConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	// 更新配置
	channel.SystemInstance.Vision.Models = form.Models

	// 保存配置
	err := channel.SystemInstance.SaveConfig()

	c.JSON(http.StatusOK, gin.H{
		"status": err == nil,
		"error":  utils.GetError(err),
	})
}

// RefreshVisionConfigAPI 刷新视觉模型配置(使配置生效)
func RefreshVisionConfigAPI(c *gin.Context) {
	// 重新加载配置到内存
	utils.RefreshCustomVisionModels(channel.SystemInstance.GetVisionModels())

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"message": "Vision models configuration refreshed successfully",
	})
}

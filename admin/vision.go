package admin

import (
	"chat/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VisionConfigForm struct {
	Models           []string `json:"models"`
	TreatAllAsVision bool     `json:"treat_all_as_vision"`
}

// GetVisionConfigAPI 获取视觉模型配置
func GetVisionConfigAPI(c *gin.Context) {
	cfg := utils.GetVisionConfig()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   cfg,
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

	// 使用独立配置保存函数
	err := utils.UpdateVisionConfig(form.Models, form.TreatAllAsVision)

	c.JSON(http.StatusOK, gin.H{
		"status": err == nil,
		"error":  utils.GetError(err),
	})
}

// RefreshVisionConfigAPI 刷新视觉模型配置(使配置生效)
func RefreshVisionConfigAPI(c *gin.Context) {
	// 重新加载配置到内存
	cfg := utils.GetVisionConfig()
	utils.RefreshCustomVisionModels(cfg.Models)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Vision models configuration refreshed successfully",
	})
}

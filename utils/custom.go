/*
 * @Author: 刘家荣 14731753+liujiarong2@user.noreply.gitee.com
 * @Date: 2025-01-21 18:04:14
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2025-10-18 01:50:00
 * @FilePath: /coai/utils/custom.go
 * @Description: 视觉模型配置管理,从系统配置读取
 */
package utils

import (
	"chat/globals"
	"fmt"
	"sync"
)

var (
	customVisionModels []string
	configMutex        sync.RWMutex
)

// RefreshCustomVisionModels 刷新自定义视觉模型列表
func RefreshCustomVisionModels(models []string) {
	configMutex.Lock()
	defer configMutex.Unlock()

	customVisionModels = make([]string, len(models))
	copy(customVisionModels, models)

	globals.Info(fmt.Sprintf("Vision models configuration refreshed, total %d models: %v", len(customVisionModels), customVisionModels))
}

// IsCustomVisionModel 检查模型是否是自定义视觉模型
func IsCustomVisionModel(model string) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()

	for _, customModel := range customVisionModels {
		if customModel == model {
			return true
		}
	}
	return false
}

// GetCustomVisionModels 获取自定义视觉模型列表
func GetCustomVisionModels() []string {
	configMutex.RLock()
	defer configMutex.RUnlock()

	result := make([]string, len(customVisionModels))
	copy(result, customVisionModels)
	return result
}

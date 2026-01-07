/*
 * @Author: Liu Jiarong
 * @Date: 2026-01-07
 * @Description: 视觉模型独立配置管理，避免被 viper.WriteConfig() 覆盖
 */
package utils

import (
	"chat/globals"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// VisionConfig 视觉模型配置结构
type VisionConfig struct {
	Models           []string `yaml:"models" json:"models"`
	TreatAllAsVision bool     `yaml:"treat_all_as_vision" json:"treat_all_as_vision"`
}

var (
	visionConfig     *VisionConfig
	visionConfigOnce sync.Once
	visionConfigMu   sync.RWMutex
)

const visionConfigFile = "config/vision.yaml"

// LoadVisionConfig 加载视觉配置
func LoadVisionConfig() *VisionConfig {
	visionConfigOnce.Do(func() {
		visionConfig = &VisionConfig{
			Models:           []string{},
			TreatAllAsVision: false,
		}
		reloadVisionConfig()
	})
	return visionConfig
}

// reloadVisionConfig 重新加载配置文件
func reloadVisionConfig() {
	visionConfigMu.Lock()
	defer visionConfigMu.Unlock()

	data, err := os.ReadFile(visionConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			globals.Info("[vision_config] config file not found, using defaults")
			// 创建默认配置文件
			SaveVisionConfigInternal(visionConfig)
		} else {
			globals.Warn(fmt.Sprintf("[vision_config] failed to read config: %s", err.Error()))
		}
		return
	}

	var cfg VisionConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		globals.Warn(fmt.Sprintf("[vision_config] failed to parse config: %s", err.Error()))
		return
	}

	visionConfig.Models = cfg.Models
	visionConfig.TreatAllAsVision = cfg.TreatAllAsVision

	globals.Info(fmt.Sprintf("[vision_config] loaded %d vision models, treat_all_as_vision=%v",
		len(visionConfig.Models), visionConfig.TreatAllAsVision))
}

// GetVisionConfig 获取视觉配置（只读）
func GetVisionConfig() VisionConfig {
	visionConfigMu.RLock()
	defer visionConfigMu.RUnlock()

	if visionConfig == nil {
		return VisionConfig{Models: []string{}, TreatAllAsVision: false}
	}

	// 返回副本
	result := VisionConfig{
		Models:           make([]string, len(visionConfig.Models)),
		TreatAllAsVision: visionConfig.TreatAllAsVision,
	}
	copy(result.Models, visionConfig.Models)
	return result
}

// UpdateVisionConfig 更新视觉配置
func UpdateVisionConfig(models []string, treatAllAsVision bool) error {
	visionConfigMu.Lock()
	defer visionConfigMu.Unlock()

	if visionConfig == nil {
		visionConfig = &VisionConfig{}
	}

	visionConfig.Models = make([]string, len(models))
	copy(visionConfig.Models, models)
	visionConfig.TreatAllAsVision = treatAllAsVision

	// 同步刷新全局视觉模型列表
	RefreshCustomVisionModels(models)
	globals.TreatAllAsVision = treatAllAsVision

	return SaveVisionConfigInternal(visionConfig)
}

// SaveVisionConfigInternal 内部保存函数（需要在锁内调用或已持有锁）
func SaveVisionConfigInternal(cfg *VisionConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal vision config: %w", err)
	}

	// 添加注释头
	header := "# 视觉模型配置\n# 此文件独立于主配置文件，不受 viper.WriteConfig() 影响\n\n"
	content := header + string(data)

	if err := os.WriteFile(visionConfigFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write vision config: %w", err)
	}

	globals.Info("[vision_config] config saved successfully")
	return nil
}

// GetVisionModels 获取视觉模型列表
func GetVisionModels() []string {
	cfg := GetVisionConfig()
	return cfg.Models
}

// IsTreatAllAsVision 是否将所有模型视为视觉模型
func IsTreatAllAsVision() bool {
	cfg := GetVisionConfig()
	return cfg.TreatAllAsVision
}

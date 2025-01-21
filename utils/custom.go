/*
 * @Author: 刘家荣 14731753+liujiarong2@user.noreply.gitee.com
 * @Date: 2025-01-21 18:04:14
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2025-01-21 22:13:33
 * @FilePath: /coai/utils/custom.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package utils

import (
	"chat/globals"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Config struct {
	CustomVisionModels []string `json:"custom_vision_models"`
}

var (
	customVisionConfig   Config
	lastConfigUpdateTime time.Time
	configMutex          sync.RWMutex
)

func init() {
	loadConfig()
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			loadConfig()
		}
	}()
}

func loadConfig() {
	file, err := os.Open("config/vision_models.json")
	if err != nil {
		globals.Warn(fmt.Sprintf("Error opening config/vision_models.json: %v", err))
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		globals.Warn(fmt.Sprintf("Error reading config/vision_models.json: %v", err))
		return
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		globals.Warn(fmt.Sprintf("Error unmarshalling config/vision_models.json: %v", err))
		return
	}

	configMutex.Lock() // 加锁
	customVisionConfig = config
	lastConfigUpdateTime = time.Now()
	configMutex.Unlock() // 解锁
	globals.Info(fmt.Sprintf("Successfully load or update config, current config:%v, update time: %v", config, lastConfigUpdateTime))
}

func IsCustomVisionModel(model string) bool {
	configMutex.RLock()
	defer configMutex.RUnlock()
	for _, customModel := range customVisionConfig.CustomVisionModels {
		if customModel == model {
			return true
		}
	}
	return false
}

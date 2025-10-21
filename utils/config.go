package utils

import (
	"chat/globals"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var configFile = "config/config.yaml"
var configExampleFile = "config.example.yaml"

var redirectRoutes = []string{
	"/v1",
	"/mj",
	"/attachments",
}

func ReadConf() {
	viper.SetConfigFile(configFile)

	if !IsFileExist(configFile) {
		fmt.Println(fmt.Sprintf("[service] config.yaml not found, creating one from template: %s", configExampleFile))
		if err := CopyFile(configExampleFile, configFile); err != nil {
			fmt.Println(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	normalizeLinuxDoConfig()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if timeout := viper.GetInt("max_timeout"); timeout > 0 {
		globals.HttpMaxTimeout = time.Second * time.Duration(timeout)
		globals.Debug(fmt.Sprintf("[service] http client timeout set to %ds from env", timeout))
	}
}

// normalizeLinuxDoConfig ensures legacy linuxdo keys without underscores map to the canonical linux_do keys.
func normalizeLinuxDoConfig() {
	const canonicalBase = "system.oauth.linux_do"
	const altBase = "system.oauth.linuxdo"

	if !viper.IsSet(altBase) {
		return
	}

	raw := viper.GetStringMap(altBase)
	if len(raw) == 0 {
		return
	}

	setBoolIfMissing(canonicalBase+".enabled", raw, "enabled")
	setStringIfMissing(canonicalBase+".client_id", raw, "client_id", "clientid")
	setStringIfMissing(canonicalBase+".client_secret", raw, "client_secret", "clientsecret")
	setStringIfMissing(canonicalBase+".redirect_url", raw, "redirect_url", "redirecturl")
}

func setBoolIfMissing(path string, source map[string]any, key string) {
	if viper.IsSet(path) {
		return
	}

	val, ok := source[key]
	if !ok {
		return
	}

	if boolVal, valid := coerceBool(val); valid {
		viper.Set(path, boolVal)
	}
}

func setStringIfMissing(path string, source map[string]any, keys ...string) {
	if viper.IsSet(path) {
		return
	}

	if strVal, ok := findStringValue(source, keys...); ok {
		viper.Set(path, strVal)
	}
}

func coerceBool(val any) (bool, bool) {
	switch v := val.(type) {
	case bool:
		return v, true
	case string:
		str := strings.TrimSpace(v)
		if str == "" {
			return false, false
		}

		parsed, err := strconv.ParseBool(str)
		if err == nil {
			return parsed, true
		}
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return num != 0, true
		}
	case int:
		return v != 0, true
	case int8:
		return v != 0, true
	case int16:
		return v != 0, true
	case int32:
		return v != 0, true
	case int64:
		return v != 0, true
	case uint:
		return v != 0, true
	case uint8:
		return v != 0, true
	case uint16:
		return v != 0, true
	case uint32:
		return v != 0, true
	case uint64:
		return v != 0, true
	case float32:
		return v != 0, true
	case float64:
		return v != 0, true
	}

	if num, err := strconv.ParseFloat(fmt.Sprint(val), 64); err == nil {
		return num != 0, true
	}

	return false, false
}

func findStringValue(source map[string]any, keys ...string) (string, bool) {
	for _, key := range keys {
		if val, ok := source[key]; ok {
			if val == nil {
				continue
			}
			switch v := val.(type) {
			case string:
				return v, true
			case fmt.Stringer:
				return v.String(), true
			case []byte:
				return string(v), true
			}

			return fmt.Sprint(val), true
		}
	}

	return "", false
}

func NewEngine() *gin.Engine {
	if viper.GetBool("debug") {
		return gin.Default()
	}

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	return engine
}

func ApplySeo(title, icon string) {
	// seo optimization

	if !viper.GetBool("serve_static") {
		return
	}

	content, err := ReadFile("./app/dist/index.html")
	if err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to read index.html: %s", err.Error()))
		return
	}

	if len(title) > 0 {
		content = strings.ReplaceAll(content, "Chat Nio", title)
		content = strings.ReplaceAll(content, "chatnio", strings.ToLower(title))
	}

	if len(icon) > 0 {
		content = strings.ReplaceAll(content, "/favicon.ico", icon)
	}

	if err := WriteFile("./app/dist/index.cache.html", content, true); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to write index.cache.html: %s", err.Error()))
	}

	globals.Info("[service] seo optimization applied to index.cache.html")
}

func ApplyPWAManifest(content string) {
	// pwa manifest rewrite (site.webmanifest -> site.cache.webmanifest)

	if !viper.GetBool("serve_static") {
		return
	}

	if len(content) == 0 {
		// read from site.webmanifest if not provided

		var err error
		content, err = ReadFile("./app/dist/site.webmanifest")
		if err != nil {
			globals.Warn(fmt.Sprintf("[service] failed to read site.webmanifest: %s", err.Error()))
			return
		}
	}

	if err := WriteFile("./app/dist/site.cache.webmanifest", content, true); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to write site.cache.webmanifest: %s", err.Error()))
	}

	globals.Info("[service] pwa manifest applied to site.cache.webmanifest")
}

func ReadPWAManifest() (content string) {
	// read site.cache.webmanifest content or site.webmanifest if not found

	if !viper.GetBool("serve_static") {
		return
	}

	if text, err := ReadFile("./app/dist/site.cache.webmanifest"); err == nil && len(text) > 0 {
		return text
	}

	if text, err := ReadFile("./app/dist/site.webmanifest"); err != nil {
		globals.Warn(fmt.Sprintf("[service] failed to read site.webmanifest: %s", err.Error()))
	} else {
		content = text
	}

	return
}

func RegisterStaticRoute(engine *gin.Engine) {
	// static files are in ~/app/dist

	if !viper.GetBool("serve_static") {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{"status": false, "message": "not found or method not allowed"})
		})
		return
	}

	if !IsFileExist("./app/dist") {
		fmt.Println("[service] app/dist not found, please run `npm run build`")
		return
	}

	ApplySeo(viper.GetString("system.general.title"), viper.GetString("system.general.logo"))
	ApplyPWAManifest(viper.GetString("system.general.pwamanifest"))

	engine.GET("/", func(c *gin.Context) {
		c.File("./app/dist/index.cache.html")
	})

	engine.GET("/site.webmanifest", func(c *gin.Context) {
		c.File("./app/dist/site.cache.webmanifest")
	})

	engine.Use(static.Serve("/", static.LocalFile("./app/dist", true)))
	engine.NoRoute(func(c *gin.Context) {
		c.File("./app/dist/index.cache.html")
	})

	for _, route := range redirectRoutes {
		engine.Any(fmt.Sprintf("%s/*path", route), func(c *gin.Context) {
			c.Request.URL.Path = "/api" + c.Request.URL.Path
			fmt.Println(c.Request.URL.Path)
			engine.HandleContext(c)
		})
	}

	fmt.Println(`[service] start serving static files from ~/app/dist`)
}

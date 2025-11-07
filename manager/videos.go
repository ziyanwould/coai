package manager

import (
	adaptercommon "chat/adapter/common"
	"chat/admin/analysis"
	"chat/auth"
	"chat/channel"
	"chat/globals"
	"chat/utils"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

func VideosRelayAPI(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			globals.Warn(fmt.Sprintf("caught panic from chat completions api: %s (client: %s)\n%s",
				err, c.ClientIP(), stack,
			))
		}
	}()

	if globals.CloseRelay {
		abortWithErrorResponse(c, fmt.Errorf("relay api is denied of access"), "access_denied_error")
		return
	}

	username := utils.GetUserFromContext(c)
	if username == "" {
		abortWithErrorResponse(c, fmt.Errorf("access denied for invalid api key"), "authentication_error")
		return
	}

	if utils.GetAgentFromContext(c) != "api" {
		abortWithErrorResponse(c, fmt.Errorf("access denied for invalid agent"), "authentication_error")
		return
	}

	var form RelayVideoForm
	if err := c.ShouldBindJSON(&form); err != nil {
		abortWithErrorResponse(c, fmt.Errorf("invalid request body: %s", err.Error()), "invalid_request_error")
		return
	}

	prompt := strings.TrimSpace(form.Prompt)
	if prompt == "" {
		sendErrorResponse(c, fmt.Errorf("prompt is required"), "invalid_request_error")
	}

	db := utils.GetDBFromContext(c)
	cache := utils.GetCacheFromContext(c)
	user := &auth.User{
		Username: username,
	}

	form.Model = strings.TrimSuffix(form.Model, "-official")

	if form.Model == "" {
		form.Model = globals.Sora2
	}

	messages := []globals.Message{
		{Role: globals.User, Content: prompt},
	}
	check, plan := checkEnableState(db, cache, user, form.Model, messages)
	if check != nil {
		sendErrorResponse(c, check, "quota_exceeded_error")
		return
	}

	buffer := utils.NewBuffer(form.Model, messages, channel.ChargeInstance.GetCharge(form.Model))
	buffer.SetTokenName(globals.ApiTokenType)

	props := adaptercommon.CreateVideoProps(&adaptercommon.VideoProps{
		Model:          form.Model,
		Prompt:         prompt,
		Seconds:        form.Seconds,
		Size:           form.Size,
		InputReference: form.InputReference,
	})
	props.User = auth.GetUsernameString(db, user)

	group := auth.GetGroup(db, user)

	var jobJson string
	hit, err := channel.NewVideoRequestWithCache(cache, buffer, group, props, func(data *globals.Chunk) error {
		if data != nil {
			jobJson = data.Content
		}
		return nil
	})

	analysis.AnalyseRequest(form.Model, props.User, buffer, err)
	if err != nil {
		auth.RevertSubscriptionUsage(db, cache, user, form.Model)
		globals.Warn(fmt.Sprintf("error from video request api: %s (instance: %s, client: %s)", err, form.Model, c.ClientIP()))
		sendErrorResponse(c, err)
		return
	}

	if !hit {
		CollectQuota(c, user, buffer, plan, err)
	}

	job, jerr := utils.UnmarshalString[RelayVideoJob](jobJson)
	if jerr != nil {
		sendErrorResponse(c, fmt.Errorf("invalid job response: %s", jerr.Error()))
		return
	}

	c.JSON(http.StatusOK, job)
}

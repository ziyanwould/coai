package manager

import (
	"chat/auth"
	"chat/manager/conversation"
	"chat/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebsocketAuthForm struct {
	Token string `json:"token" binding:"required"`
	Id    int64  `json:"id" binding:"required"`
	Ref   string `json:"ref"`
}

func ParseAuth(c *gin.Context, token string) *auth.User {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}

	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}

	if strings.HasPrefix(token, "sk-") {
		return auth.ParseApiKey(c, token)
	}

	return auth.ParseToken(c, token)
}

func splitMessage(message string) (int, string, error) {
	parts := strings.SplitN(message, ":", 2)
	if len(parts) == 2 {
		if id, err := strconv.Atoi(parts[0]); err == nil {
			return id, parts[1], nil
		}
	}

	return 0, message, fmt.Errorf("message type error")
}

func getId(message string) (int, error) {
	if id, err := strconv.Atoi(message); err == nil {
		return id, nil
	}

	return 0, fmt.Errorf("message type error")
}

func ChatAPI(c *gin.Context) {
	var conn *utils.WebSocket
	if conn = utils.NewWebsocket(c, false); conn == nil {
		return
	}
	defer conn.DeferClose()

	db := utils.GetDBFromContext(c)

	form, err := utils.ReadForm[WebsocketAuthForm](conn)
	if err != nil {
		return
	}

	user := ParseAuth(c, form.Token)
	authenticated := user != nil

	id := auth.GetId(db, user)

	instance := conversation.ExtractConversation(db, user, form.Id, form.Ref)
	hash := fmt.Sprintf(":chatthread:%s", utils.Md5Encrypt(utils.Multi(
		authenticated,
		strconv.FormatInt(id, 10),
		c.ClientIP(),
	)))

	buf := NewConnection(conn, authenticated, hash, 10)
	buf.Handle(func(form *conversation.FormMessage) error {
		switch form.Type {
		case ChatType:
			if instance.HandleMessage(db, form) {
				response := ChatHandler(buf, user, instance, false, getClientIP(c))
				instance.SaveResponse(db, response)
			}
		case StopType:
			break
		case ShareType:
			instance.LoadSharing(db, form.Message)
		case RestartType:
			// reset the params if set
			instance.ApplyParam(form)

			response := ChatHandler(buf, user, instance, true, getClientIP(c))
			instance.SaveResponse(db, response)
		case MaskType:
			instance.LoadMask(form.Message)
		case EditType:
			if id, message, err := splitMessage(form.Message); err == nil {
				instance.EditMessage(id, message)
				instance.SaveConversation(db)
			} else {
				return err
			}
		case RemoveType:
			id, err := getId(form.Message)
			if err != nil {
				return err
			}

			instance.RemoveMessage(id)
			instance.SaveConversation(db)
		}

		return nil
	})
}
func getClientIP(c *gin.Context) string {
	// 尝试从 X-Forwarded-For 头中获取 IP 地址
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		addresses := strings.Split(forwardedFor, ",")
		// 获取第一个非空的 IP 地址
		for _, addr := range addresses {
			if ip := strings.TrimSpace(addr); ip != "" {
				return ip
			}
		}
	}

	// 尝试从 X-Real-IP 头中获取 IP 地址
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 如果以上方法都失败，则使用 c.ClientIP()
	return c.ClientIP()
}

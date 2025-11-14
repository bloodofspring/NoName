package commandblock

import (
	"app/internal/handlers"
	"app/internal/handlers/shared"
	"time"

	"app/pkg/database/models"
	e "app/pkg/errors"

	"github.com/go-pg/pg/v10"
	tele "gopkg.in/telebot.v4"
)

func CommandBlockUserChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(
		10*time.Second,
		shared.ConnectDatabase,
		GetUserFromThreadID,
		SuccessMessage,
	)
}

func GetUserFromThreadID(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := (*args)["db"].(*pg.DB)

	if c.Chat().Type != tele.ChatSuperGroup {
		c.Reply("Chat should be a supergroup")
		return args, e.Nil()
	}

	threadId := c.Message().ThreadID

	var thread models.Thread
	err := db.Model(&thread).
		Where("thread_id = ?", threadId).
		Where("chat_id = ?", c.Chat().ID).
		Select()
	if err != nil {
		return args, e.Error(err, "Failed to select thread").WithSeverity(e.Critical).WithData(map[string]any{
			"thread_id": threadId,
		})
	}

	_, err = db.Model(&models.User{
		TgID: thread.AssociatedUserID,
	}).
		WherePK().
		Set("is_blocked = true").
		Update()
	if err != nil {
		return args, e.Error(err, "Failed to update user").WithSeverity(e.Critical).WithData(map[string]any{
			"user": thread.AssociatedUserID,
		})
	}

	return args, e.Nil()
}

func SuccessMessage(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	c.Reply("User blocked successfully")

	return args, e.Nil()
}
package commandblock

import (
	"app/internal/handlers"
	"time"

	"app/pkg/database"
	"app/pkg/database/models"
	e "app/pkg/errors"

	tele "gopkg.in/telebot.v4"
)

func CommandBlockUserChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(
		10*time.Second,
		handlers.InitChainHandler(GetUserFromThreadID),
		handlers.InitChainHandler(SuccessMessage),
	)
}

func GetUserFromThreadID(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := database.GetDB()

	if c.Chat().Type != tele.ChatSuperGroup {
		c.Reply("Chat should be a supergroup")
		return args, e.NewError("chat should be a supergroup", "Chat should be a supergroup").WithSeverity(e.Ingnored).WithData(map[string]any{
			"sender": (*args)["sender"],
		})
	}

	threadId := c.Message().ThreadID

	if threadId == 0 {
		c.Reply("Thread not found")
		return args, e.NewError("thread not found", "Thread not found").WithSeverity(e.Notice).WithData(map[string]any{
			"chat_id": c.Chat().ID,
		})
	}

	var thread models.Thread
	err := db.Model(&thread).
		Where("thread_id = ?", threadId).
		Where("chat_id = ?", c.Chat().ID).
		Select()
	if err != nil {
		return args, e.FromError(err, "Failed to select thread").WithSeverity(e.Notice).WithData(map[string]any{
			"thread_id": threadId,
			"chat_id": c.Chat().ID,
		})
	}

	_, err = db.Model(&models.User{
		TgID: thread.AssociatedUserID,
	}).
		WherePK().
		Set("is_blocked = true").
		Update()
	if err != nil {
		return args, e.FromError(err, "Failed to update user").WithSeverity(e.Notice).WithData(map[string]any{
			"user": thread.AssociatedUserID,
		})
	}

	return args, e.Nil()
}

func SuccessMessage(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	c.Reply("User blocked successfully")

	return args, e.Nil()
}
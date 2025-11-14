package commandaddchat

import (
	"app/internal/handlers"
	"app/internal/handlers/shared"
	"fmt"
	"time"

	"app/pkg/database/models"
	e "app/pkg/errors"

	"github.com/go-pg/pg/v10"
	tele "gopkg.in/telebot.v4"
)

func CommandAddChatChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(
		10*time.Second,
		shared.ConnectDatabase,
		GetSenderAndTargetUser,
		InitNewChatForUser,
		SendChatAddedMessage,
	)
}

func GetSenderAndTargetUser(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := (*args)["db"].(*pg.DB)
	sender := &models.User{
		TgID: c.Sender().ID,
	}

	err := db.Model(sender).WherePK().Select()
	if err != nil {
		return args, e.Error(err, "Failed to select sender user").WithSeverity(e.Critical).WithData(map[string]any{
			"sender": sender,
		})
	}

	(*args)["sender"] = sender

	return args, e.Nil()
}

func InitNewChatForUser(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := (*args)["db"].(*pg.DB)

	if c.Chat().Type != tele.ChatSuperGroup {
		c.Reply("Chat should be a supergroup")
		return args, e.Nil()
	}
	
	var chat models.Chat
	err := db.Model(&chat).
		Where("chat_owner_id = ?", (*args)["sender"].(*models.User).TgID).
		Select()

	if err == nil {
		c.Reply("Chat already exists. Replace with new one? (TODO: Implement)")
		return args, e.Nil()
	}
	
	if err != pg.ErrNoRows {
		return args, e.Error(err, "Failed to select chat").WithSeverity(e.Critical).WithData(map[string]any{
			"sender": (*args)["sender"],
		})
	}

	chat = models.Chat{
		ChatOwnerID: (*args)["sender"].(*models.User).TgID,
		TgID: c.Chat().ID,
	}
	_, err = db.Model(&chat).Insert()
	if err != nil {
		return args, e.Error(err, "Failed to insert chat").WithSeverity(e.Critical).WithData(map[string]any{
			"chat": chat,
		})
	}

	return args, e.Nil()
}

func SendChatAddedMessage(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	c.Reply(fmt.Sprintf("Chat %s added successfully", c.Chat().Title))
	return args, e.Nil()
}
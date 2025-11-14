package newmessage

import (
	"app/internal/handlers"
	"app/internal/handlers/shared"
	"time"

	tele "gopkg.in/telebot.v4"
	e "app/pkg/errors"
	"app/pkg/database/models"
)

func NewMessageChain() *handlers.HandlerChain {
	return handlers.HandlerChain{}.Init(
		10*time.Second,
		shared.ConnectDatabase,
		shared.GetSenderAndTargetUser,
		shared.GetOrCrateThread,
		ChackUserIsBlocked,
		RedirectMessageToThread,
	)
}

func ChackUserIsBlocked(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	user := (*args)["user"].(*models.User)
	if user.IsBlocked {
		c.Reply("You are blocked. There is nothting you can do with it :3")
		return args, e.Nil()
	}

	return args, e.Nil()
}

func RedirectMessageToThread(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	thread := (*args)["thread"].(*models.Thread)

	chatRecipient := &tele.Chat{ID: thread.ChatID, Type: tele.ChatSuperGroup}
	options := &tele.SendOptions{ThreadID: thread.ThreadID}

	c.ForwardTo(
		chatRecipient,
		options,
	)

	return args, e.Nil()
}

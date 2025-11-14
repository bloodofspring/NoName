package shared

// TODO: Внести возможность создавать список зависимостей для shared функций

import (
	"app/internal/handlers"
	"app/pkg/database"
	"app/pkg/database/models"
	e "app/pkg/errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	tele "gopkg.in/telebot.v4"
)

func ConnectDatabase(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := database.GetDB()

	newArgs := make(handlers.Arg)
	newArgs["db"] = db

	return &newArgs, e.Nil()
}

func GetSenderAndTargetUser(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := (*args)["db"].(*pg.DB)
	user := &models.User{
		TgID: c.Sender().ID,
	}

	target := &models.User{
		TgID: viper.GetInt64("OWNER_TG_ID"),
	}

	err := db.Model(target).WherePK().Select()
	if err != nil {
		return args, e.Error(err, "Failed to select target user").WithSeverity(e.Critical).WithData(map[string]any{
			"target": target,
		})
	}

	_, err = db.Model(user).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return args, e.Error(err, "Failed to insert user").WithSeverity(e.Critical)
	}

	(*args)["user"] = user
	(*args)["target"] = target

	return args, e.Nil()
}

func GetOrCrateThread(c tele.Context, args *handlers.Arg) (*handlers.Arg, *e.ErrorInfo) {
	db := (*args)["db"].(*pg.DB)

	var Chat *models.Chat
	err := db.Model(&Chat).
		Where("chat_owner_id = ?", (*args)["target"].(*models.User).TgID).
		Select()
	if err != nil {
		return args, e.Error(err, "Failed to select chat").WithSeverity(e.Critical).WithData(map[string]any{
			"target": (*args)["target"],
		})
	}

	(*args)["chat"] = Chat

	var thread models.Thread

	err = db.Model(&thread).
		Where("chat_id = ?", Chat.TgID).
		Where("associated_user_id = ?", (*args)["user"].(*models.User).TgID).
		Select()

	if err == nil {
		(*args)["thread"] = thread
		return args, e.Nil()
	}

	if err != pg.ErrNoRows {
		return args, e.Error(err, "Failed to select thread").WithSeverity(e.Critical).WithData(map[string]any{
			"target": (*args)["target"],
			"user": (*args)["user"],
		})
	}

	// TODO: Создать массив с IDшниками иконок и выбирать случайную
	t, err := c.Bot().CreateTopic(
		&tele.Chat{
			ID: Chat.TgID,
		},
		&tele.Topic{
			Name: fmt.Sprintf("@%s", c.Sender().Username),
			IconCustomEmojiID: "5199590728270886590",
		},
	)

	if err != nil {
		return args, e.Error(err, "Failed to create topic").WithSeverity(e.Critical).WithData(map[string]any{
			"chat": Chat,
			"user": (*args)["user"],
		})
	}

	thread = models.Thread{
		ThreadID: t.ThreadID,
		ChatID: Chat.TgID,
		AssociatedUserID: (*args)["user"].(*models.User).TgID,
	}

	_, err = db.Model(&thread).Insert()
	if err != nil {
		return args, e.Error(err, "Failed to insert thread").WithSeverity(e.Critical).WithData(map[string]any{
			"thread": thread,
			"chat": Chat,
			"user": (*args)["user"],
		})
	}

	(*args)["thread"] = thread

	return args, e.Nil()
}
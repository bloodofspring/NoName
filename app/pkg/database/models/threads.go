package models

type Thread struct {
	ID int `pg:",pk,notnull,autoincrement"`
	ThreadID int64 `pg:",notnull"`
	CreatedAt int64 `pg:"default:extract(epoch from now())"`
	UpdatedAt int64 `pg:"default:extract(epoch from now())"`

	ChatID int64 `pg:",notnull"`
	Chat *Chat `pg:"rel:has-one,fk:chat_id"`
	
	AssociatedUserID int64 `pg:",notnull"`
	AssociatedUser *User `pg:"rel:has-one,fk:associated_user_id"`
}
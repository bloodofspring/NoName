package models

type ResentMessage struct {
	ID        int   `pg:"id,pk,notnull"`
	CreatedAt int64 `pg:"default:extract(epoch from now())"`
	UpdatedAt int64 `pg:"default:extract(epoch from now())"`

	ThreadID int     `pg:",notnull"`
	Thread   *Thread `pg:"rel:has-one,fk:thread_id"`

	ChatID int64 `pg:",notnull"`
	Chat   *Chat `pg:"rel:has-one,fk:chat_id"`

	SenderChatMessageID int `pg:",notnull"`
	TargetChatMessageID int `pg:",notnull"`
}

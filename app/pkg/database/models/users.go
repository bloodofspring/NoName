package models

type User struct {
	TgID int64 `pg:",pk,notnull"`
	CreatedAt int64 `pg:"default:extract(epoch from now())"`
	UpdatedAt int64 `pg:"default:extract(epoch from now())"`
	
	IsOwner bool `pg:"default:false"`
	IsBlocked bool `pg:"default:false"`

	OwnedChats []*Chat `pg:"rel:has-many,fk:chat_owner_id"`
	AssociatedThreads []*Thread `pg:"rel:has-many,fk:associated_user_id"`
}

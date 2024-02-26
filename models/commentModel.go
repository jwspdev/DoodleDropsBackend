package models

type Comment struct {
	AuthorID uint64
	Author   User
}

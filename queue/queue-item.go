package queue

import "github.com/google/uuid"

type KokaqItem struct {
	Id       uuid.UUID
	Priority int
}

package auth

import "sync"

const constantUserID = 1

var (
	once          sync.Once
	currentUserID uint
)

func GetCurrentUserID() uint {
	once.Do(func() {
		currentUserID = constantUserID
	})

	return currentUserID
}

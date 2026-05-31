package utils

import "fmt"

func SessionKey(id uint) string {
	return fmt.Sprintf("session:%s", id)
}

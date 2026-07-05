package utils

import "fmt"

func SessionKey(id uint) string {
	return fmt.Sprintf("session:%v", id)
}

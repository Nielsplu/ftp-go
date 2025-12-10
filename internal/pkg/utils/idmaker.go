package utils

import (
	"fmt"
	"time"
)

func MakeId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
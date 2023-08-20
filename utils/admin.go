package utils

import (
	"github.com/skifli/creft/database"
)

func HasAdminPerms(userID string) bool {
	if _, ok := database.DatabaseJSON["admins"].(map[string]any)[userID]; ok {
		return true
	}

	return false
}

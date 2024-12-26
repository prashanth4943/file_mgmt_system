package helper

import (
	"context"
	"file_mgmt_system/middleware"
)

// Helper function to extract the email from the context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(middleware.UserKey).(string)
	return email, ok
}

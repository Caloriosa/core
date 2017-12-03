package sanitization

import "core/types"

func UserSanitization(user *types.User) {
	user.Password = ""
}

package sanitization

import "core/types"

func UserSanitization(user *types.User, strict bool) {
	user.Password = ""
	if strict {
		user.Email = ""
	}
}

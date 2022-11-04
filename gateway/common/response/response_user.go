package response

import "gateway/model"

type UserResponse struct {
	Token string
	User  model.User
}

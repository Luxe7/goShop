package response

import "time"

type UserResponse struct {
	Id       uint32    `json:"id"`
	NickName string    `json:"name"`
	Birthday time.Time `json:"birthday"`
	Gender   string    `json:"gender"`
	Mobile   string    `json:"mobile"`
}

package forms

type PassWordLoginForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	PassWord string `json:"password" binding:"required,min=3,max=20"`
}

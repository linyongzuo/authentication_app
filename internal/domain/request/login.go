package request

type AdminLoginReq struct {
	Header
	UserName string `json:"userName"` // 用户名
	Password string `json:"password"` // 密码
}

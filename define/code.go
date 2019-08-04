package define


type ErrorCode int16

// 错误提示
const (
	Success       ErrorCode = iota
	RoomNonExist
	CardNonEnough
	UserOffline
	UserNonExist
	ServerError
	NodeError
	FormatError
	TokenInvalid
	RepeatLogin
)

var ErrorCodeDesc = [...]string{
	Success:       "OK",
	RoomNonExist:  "房间不存在",
	CardNonEnough: "房卡不足",
	UserOffline:   "用户离线",
	UserNonExist:  "用户不存在",
	ServerError:   "服务器错误",
	FormatError:   "协议格式错误",
	TokenInvalid:  "无效身份",
	RepeatLogin:   "账号重复登录",
	NodeError:     "节点连接错误",
}



func (s ErrorCode) String() string {
	return ErrorCodeDesc[s]
}

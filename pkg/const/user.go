package _const

const CurrentUser = "user"

type Requester interface {
	GetRole() string
	GetUserName() string
	GetId() string
}

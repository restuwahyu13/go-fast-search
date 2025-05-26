package exception

import (
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type usersException struct{}

func NewUsersException() inf.IUsersException {
	return usersException{}
}

func (e usersException) CreateUsers(key string) string {
	msg := make(map[string]string)

	msg["users_exists"] = "User already exists in our system"
	msg["create_users_failed"] = "Failed to create new users"

	return msg[key]
}

func (e usersException) UpdateUsers(key string) string {
	msg := make(map[string]string)

	msg["users_notfoud"] = "User is not exists in our system"
	msg["update_users_failed"] = "Failed to update new users"

	return msg[key]
}

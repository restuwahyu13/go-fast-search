package exception

import (
	inf "github.com/restuwahyu13/go-fast-search/shared/interfaces"
)

type usersException struct{}

func NewUsersException() inf.IUsersException {
	return usersException{}
}

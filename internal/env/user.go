package env

import "os/user"

var u *user.User

func getUser() *user.User {
	if u != nil {
		return u
	}
	var err error
	u, err = user.Current()
	if err != nil {
		panic(err)
	}
	return u
}

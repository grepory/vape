package servicer

import (
	"github.com/opsee/vape/model"
)

func SendTemplatedEmail(userId int, template string, vars map[string]string) (*model.User, error) {
	if vars == nil {
		vars = make(map[string]string)
	}

	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	go func() {
		mailTemplatedMessage(user.Email, user.Name, template, vars)
	}()

	return user, nil
}

package dbUser

import (
	"fmt"
)

func (dbObj *AccountType) CheckPassword(Password string) (resErr error) {
	resErr = nil

	AccountData := dbObj.AccountData

	if AccountData.Password != Password {
		dbObj.DB.Close()
		return fmt.Errorf("密码错误")
	}

	return
}

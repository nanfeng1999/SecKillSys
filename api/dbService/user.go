package dbService

import (
	"SecKillSys/data"
	"SecKillSys/model"
)

func GetUserByUsername(username string)(model.User, error){
	user := model.User{}
	operation := data.Db.Where("username = ?", username).First(&user)
	return user, operation.Error
}

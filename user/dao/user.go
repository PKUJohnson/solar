package dao

import std "github.com/PKUJohnson/solar/std"
import (
	"fmt"
	"github.com/PKUJohnson/solar/user/models"
)


func CreateUser(tuuser *models.User) *models.User {
	if err := models.DB().Save(tuuser).Error; err != nil {
		std.LogErrorc("mysql", err, "failed to create tu_user")
	}

	tuuser_create := new(models.User)
	tuuser_create.UserId = tuuser.UserId // return the new userid

	return tuuser
}

func GetFUserByMobile(mobile string) *models.FUser {
	result := new(models.FUser)
	if db := models.DB().Where("mobile = ? and valid = 1", mobile).First(&result); db.RecordNotFound() {
		std.LogInfoc("mysql", fmt.Sprintf("Can't find user %s", mobile))
		return nil
	}
	return result
}

func GetUserById(user_id int64) *models.User {
	result := new(models.User)
	if db := models.DB().Where(user_id).First(&result); db.RecordNotFound() {
		std.LogInfoc("mysql", fmt.Sprintf("Can't find user %d", user_id))
		return nil
	}
	return result
}

func GetUserByFUnionId(f_union_id string) *models.User {
	result := new(models.User)
	if db := models.DB().Where("f_union_id = ? and app_id = 'solar'", f_union_id).First(&result); db.RecordNotFound() {
		std.LogInfoc("mysql", fmt.Sprintf("Can't find user %s", f_union_id))
		return nil
	}
	return result
}
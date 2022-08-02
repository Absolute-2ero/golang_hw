package model

import "gorm.io/gorm"

//hint: 如果你想直接返回结构体，可以考虑在这里加上`json`的tag
type User struct {
	ID       uint   `gorm:"not null;autoIncrement" json:"id"`          //用户id
	UserName string `gorm:"type:varchar(255);unique;" json:"username"` //用户名字
	Votes    uint   `gorm:"default:0" json:"votes"`                    //用户所得票数
}

//TODO: 添加相应的与数据库交互逻辑，这里为你添加了部分逻辑，你可以自由选用
func CreateUser(name string) (error, User) {
	user := User{
		UserName: name,
		Votes:    0,
	}
	tx := DB.Create(&user)
	return tx.Error, user
}

func GetUserByName(name string) (error, User) {
	user := User{}
	tx := DB.Model(&User{}).Where("user_name=?", name).First(&user)
	return tx.Error, user
}

func AddVoteForUser(name string) error {
	err, user := GetUserByName(name)
	if err != nil {
		return err
	}
	tx := DB.Model(&User{}).Where("id = ?", user.ID).Update("votes", gorm.Expr("votes+ ?", 1))
	return tx.Error
}

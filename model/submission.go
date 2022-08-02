package model

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

//hint: 如果你想直接返回结构体，可以考虑在这里加上`json`的tag
type Submission struct {
	ID        uint   `gorm:"not null;autoIncrement"`
	UserName  string `gorm:"type:varchar(255);"`
	Avatar    string //头像base64，也可以是一个头像链接
	CreatedAt int64  //提交时间
	Score     int    //评测成绩
	Subscore1 int    //评测小分
	Subscore2 int    //评测小分
	Subscore3 int    //评测小分
}

//这里提供返回的submission的示例结构
type ReturnSub struct {
	UserName  string `json:"user"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"time"`
	Score     int    `json:"score"`
	UserVotes int    `json:"votes"`
	Subscore1 int
	Subscore2 int
	Subscore3 int
}

type Eval struct {
	Mainscore float64
	Subscore1 float64
	Subscore2 float64
	Subscore3 float64
}

/*TODO: 添加相应的与数据库交互逻辑，补全参数和返回值，可以参考user.go的设计思路*/

func Time2Timestamp(datetime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	rt := strings.Index("/", datetime)
	if rt == -1 {
		tt, _ := time.ParseInLocation("2006/01/02 15:04:05", datetime, loc)
		return tt.Unix()
	}
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", datetime, loc)
	return tt.UnixNano()
}

func Judge(content string) Eval {
	var ev Eval
	data, err_read := ioutil.ReadFile("ground_truth.txt")
	if err_read != nil {
		fmt.Println("Failed to read file!")
		return ev
	}
	submitLines := strings.Split(content, "\n")
	var submit1, submit2, submit3 []bool
	for _, line := range submitLines {
		temp1 := strings.Split(line, ",")[1]
		temp2 := strings.Split(line, ",")[2]
		temp3 := strings.Split(line, ",")[3]
		if temp1 == "1" {
			submit1 = append(submit1, true)
		} else {
			submit1 = append(submit1, false)
		}
		if temp2 == "1" {
			submit2 = append(submit2, true)
		} else {
			submit2 = append(submit2, false)
		}
		if temp3 == "1" {
			submit3 = append(submit3, true)
		} else {
			submit3 = append(submit3, false)
		}
	}

	var ans1, ans2, ans3 []bool
	dataLines := strings.Split(string(data), "\n")
	for _, line := range dataLines {
		temp1 := strings.Split(line, ",")[1]
		temp2 := strings.Split(line, ",")[2]
		temp3 := strings.Split(line, ",")[3]
		if temp1 == "1" {
			ans1 = append(ans1, true)
		} else {
			ans1 = append(ans1, false)
		}
		if temp2 == "1" {
			ans2 = append(ans2, true)
		} else {
			ans2 = append(ans2, false)
		}
		if temp3 == "1" {
			ans3 = append(ans3, true)
		} else {
			ans3 = append(ans3, false)
		}
	}
	eval1, eval2, eval3 := 0.0, 0.0, 0.0
	var k = float64(len(ans1))
	var l = float64(len(submit1))
	if k != l {
		fmt.Println("Illegal submission")
		return ev
	}
	var i = 0
	for i < int(k) {
		if ans1[i] == submit1[i] {
			eval1 = eval1 + 1
		}
		if ans2[i] == submit2[i] {
			eval2 = eval2 + 1
		}
		if ans3[i] == submit3[i] {
			eval3 = eval3 + 1
		}
		i = i + 1
	}
	ev.Subscore1 = eval1 / l * 100
	ev.Subscore2 = eval2 / l * 100
	ev.Subscore3 = eval3 / l * 100
	ev.Mainscore = ev.Subscore1*11 + ev.Subscore2*45 + ev.Subscore3*14
	return ev
}

func CreateSubmission(username, avatar, submitted_content string) error {
	var new_sub Submission
	ev := Judge(submitted_content)
	new_sub.Score = int(ev.Mainscore)
	new_sub.Subscore1 = int(ev.Subscore1)
	new_sub.Subscore2 = int(ev.Subscore2)
	new_sub.Subscore3 = int(ev.Subscore3)
	new_sub.CreatedAt = Time2Timestamp(time.Now().String())
	if len(username) > 255 {
		new_sub.UserName = ""
	} else {
		new_sub.UserName = username
	}
	err, name := GetUserByName(username)
	if err != nil {
		temp_err, new_user := CreateUser(username)
		if temp_err != nil {
			new_user.Votes = 0
		}
	}
	name.Votes++
	name.Votes--
	new_sub.Avatar = avatar
	tx := DB.Create(&new_sub)
	return tx.Error
}

func GetUserSubmissions(username string) []ReturnSub {
	//返回某一用户的所有提交
	//在查询时可以使用.Order()来控制结果的顺序，详见https://gorm.io/zh_CN/docs/query.html#Order
	//当然，也可以查询后在这个函数里手动完成排序
	var AllSub []ReturnSub
	DB.Model(&Submission{}).Where("user" == username).Find(&AllSub).Order("time")
	return AllSub
}

func GetLeaderBoard() []ReturnSub {
	//一个可行的思路，先全部选出submission，然后手动选出每个用户的最后一次提交
	var AllSub []ReturnSub
	var AllUser []User
	DB.Model(&User{}).Where("1=1").Find(&AllUser)
	var RetList []ReturnSub
	for _, user := range AllUser {
		DB.Model(&Submission{}).Where("user" == user.UserName).Find(&AllSub).Order("score")
		RetList = append(RetList, AllSub[0])
	}
	return RetList
}

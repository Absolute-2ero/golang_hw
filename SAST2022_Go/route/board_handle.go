package route

import (
	"github.com/gin-gonic/gin"
	"leadboard/model"
	"net/http"
)

//TODO:在这里完成handle function，返回所有的leader board内容
func HandleGetBoard(g *gin.Context) {
	RetList := model.GetLeaderBoard()
	g.JSON(http.StatusAccepted, gin.H{
		"msg":         "success",
		"leaderboard": RetList,
	})
}

//TODO:在这里完成返回一个用户提交历史的Handle function
func HandleUserHistory(g *gin.Context) {
	type CreateUserForm struct {
		UserName string `json:"username"`
	}
	var form CreateUserForm
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Form",
		})
	} else if form.UserName == "" {
		g.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Username",
		})
	} else {
		RetList := model.GetUserSubmissions(form.UserName)
		g.JSON(http.StatusAccepted, gin.H{
			"msg":     "Success",
			"history": RetList,
		})
	}
}

func HandleSubmit(g *gin.Context) {
	type SubmissionBody struct {
		User    string `json:"user"`
		Content string `json:"content"`
		Avatar  string `json:"avatar"`
	}
	var submission SubmissionBody
	err := g.ShouldBindJSON(&submission)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"code": 114514, "msg": "这不是 JSON 啊"})
		return
	}

	if submission.User == "" || submission.Content == "" || submission.Avatar == "" {
		g.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "你是一个，一个一个一个不全的参数啊啊啊"})
		return
	}

	name, avatar, content := submission.User, submission.Avatar, submission.Content

	if len(name) > 255 {
		g.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "你是一个，一个一个一个太长的用户名啊啊啊"})
		return
	}

	if len(avatar) > 100000 {
		g.JSON(http.StatusBadRequest, gin.H{"code": -2, "msg": "你是一个，一个一个一个太大的图像啊啊啊"})
		return
	}

	ev := model.Judge(content)

	if ev.Mainscore == 0 {
		g.JSON(http.StatusBadRequest, gin.H{"code": -3, "msg": "你是一个，一个一个一个非法的提交内容啊啊啊"})
		return
	}

	var cnt int64
	if model.DB.Model(&model.User{}).Where("user_name = ?", name).Count(&cnt); cnt == 0 {
		_, _ = model.CreateUser(name)
	}

	_ = model.CreateSubmission(name, avatar, content)

	g.JSON(http.StatusOK, gin.H{"code": 0, "leaderboard": model.GetLeaderBoard()})
}

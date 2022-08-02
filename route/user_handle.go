package route

import (
	"github.com/gin-gonic/gin"
	"leadboard/model"
	"net/http"
)

func HandleCreateUser(g *gin.Context) {
	type CreateUserForm struct {
		UserName string `json:"user_name"`
	}
	var form CreateUserForm
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid Form",
		})
	}
	err, id := model.CreateUser(form.UserName)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Error occurred",
		})
	} else {
		g.JSON(http.StatusAccepted, gin.H{
			"msg":     "Success",
			"user_id": id,
		})
	}
}

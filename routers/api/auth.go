package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/dayu26/crucian/pkg/app"
	"github.com/dayu26/crucian/pkg/e"
	"github.com/dayu26/crucian/pkg/util"
	"github.com/dayu26/crucian/service/auth_service"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

//GetAuth function
func GetAuth(c *gin.Context) {
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		app.JsonError(c, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		app.JsonError(c, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		app.JsonError(c, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		app.JsonError(c, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	app.JsonSuccess(c, e.SUCCESS, gin.H{
		"token": token,
	})
}

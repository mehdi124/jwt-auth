package controllers


import (
	"jwt-auth/models"
	Response "jwt-auth/utils/response"
	"jwt-auth/utils/token"
	"net/http"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CurrentUser(c *gin.Context){

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		data := make(map[string]string)
		data["error"] =  err.Error()
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	u,err := models.GetUserByID(user_id)

	if err != nil {
		data := make(map[string]string)
		data["error"] =  err.Error()
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	//TODO json response fix bug
	c.JSON(http.StatusOK, gin.H{"message":"success","data":u})
}

func Login(c *gin.Context){

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		data := make(map[string]string)
		data["error"] =  err.Error()
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	token, err := models.LoginCheck(input.Username, input.Password)

	if err != nil {
		data := make(map[string]string)
		data["error"] =  "username or password is incorrect."
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	data := make(map[string]string)
	data["token"] = token
	Response.SuccessResponse(c,http.StatusOK,data)

}


type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context){

	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {

		data := make(map[string]string)
		data["error"] =  err.Error()
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	_,err := u.SaveUser()

	if err != nil{
		data := make(map[string]string)
		data["error"] =  err.Error()
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		data := make(map[string]string)
		data["error"] =  "username or password is incorrect."
		Response.ErrorResponse(c,http.StatusBadRequest,data)
		return
	}

	data := make(map[string]string)
	data["token"] = token
	Response.SuccessResponse(c,http.StatusOK,data)

}


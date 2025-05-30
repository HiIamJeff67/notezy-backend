package controllers

import (
	models "notezy-backend/app/models"

	"github.com/gin-gonic/gin"
)

/* ============================== DTO ============================== */
type RegisterDto struct {
	InputData models.CreateUserInput
}

/* ============================== DTO ============================== */

/* ============================== Controller ============================== */
func Register(ctx *gin.Context) {
	// var dto RegisterDto
	// if err := ctx.ShouldBindJSON(&dto); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// newUser, err := models.CreateUser(dto.InputData)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "success",
	// 	"data": gin.H{
	// 		"createdAt": newUser.CreatedAt,
	// 	},
	// })
}

// func Login(ctx *gin.Context) {
// 	var dto models.CreateUserDto
// 	if err := ctx.ShouldBindJSON(&dto); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
// 		return
// 	}

// }
/* ============================== Controller ============================== */

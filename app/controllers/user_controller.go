package controllers

import (
	"net/http"
	"notezy-backend/app/models/inputs"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

/* ============================== DTO ============================== */
type UpdateUserDto struct {
	Id        uuid.UUID
	InputData inputs.UpdateUserInput
}

/* ============================== DTO ============================== */

/* ============================== Controller ============================== */
// func GetAllUsers(ctx *gin.Context) {
// 	users, err := models.GetAllUsers()
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "success",
// 		"data":    users,
// 	})
// }

// func UpdateUserById(ctx *gin.Context) {
// 	var dto UpdateUserDto
// 	if err := ctx.ShouldBindJSON(&dto); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedUser, err := models.UpdateUserById(dto.Id, dto.InputData)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message":   "success",
// 		"updatedAt": updatedUser.CreatedAt,
// 	})
// }

func FindAllUsers(ctx *gin.Context) {
	resDto, exception := services.FindAllUsers()
	if exception != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": exception.Log().Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"users": resDto,
		},
	})
}

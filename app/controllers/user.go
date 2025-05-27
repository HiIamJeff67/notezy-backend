package controllers

import (
	models "go-gorm-api/app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

/* ============================== DTO ============================== */
type UpdateUserDto struct {
	ID uuid.UUID
	InputData models.UpdateUserInput
}
/* ============================== DTO ============================== */

/* ============================== Controller ============================== */
func GetAllUsers(ctx *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success", 
		"data": users,
	})
}

func UpdateUserById(ctx *gin.Context) {
	var dto UpdateUserDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		return
	}

	updatedUser, err := models.UpdateUserById(dto.ID, dto.InputData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success", 
		"updatedAt": updatedUser.CreatedAt,
	})
}
/* ============================== Controller ============================== */
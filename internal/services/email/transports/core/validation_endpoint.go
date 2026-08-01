package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	corecontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	emaildto "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

func (t Endpoint) SendValidationEmail(ctx *gin.Context) {
	request := &corecontract.Request[emaildto.SendValidationEmailRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:  corecontract.Version,
			Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exceptions.New(
				"InvalidRequest",
				"Email",
				"SendValidationEmail",
				"The email request is invalid",
				http.StatusBadRequest,
			).WithOrigin(err),
		})
		return
	}
	if request.Operation != emaildto.SendValidationEmailOperation {
		ctx.JSON(http.StatusBadRequest, corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:  corecontract.Version,
			Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exceptions.New(
				"InvalidOperation",
				"Email",
				"SendValidationEmail",
				"The email operation is invalid",
				http.StatusBadRequest,
			),
		})
		return
	}
	if err := requestValidator.Struct(request.Dto); err != nil {
		ctx.JSON(http.StatusBadRequest, corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:  corecontract.Version,
			Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exceptions.New(
				"InvalidRequest",
				"Email",
				"SendValidationEmail",
				"The email request is invalid",
				http.StatusBadRequest,
			).WithOrigin(err),
		})
		return
	}

	exception := t.sender.SendValidationEmail(
		request.Dto.To,
		request.Dto.UserName,
		request.Dto.AuthCode,
		request.Dto.UserAgent,
		request.Dto.ExpiredAt,
	)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode(), corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:   corecontract.Version,
			Metadata:  corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, corecontract.Response[emaildto.SendEmailResponseDto]{
		Version:  corecontract.Version,
		Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     emaildto.SendEmailResponseDto{QueuedAt: time.Now()},
	})
}

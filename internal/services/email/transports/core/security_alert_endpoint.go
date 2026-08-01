package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	corecontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	emaildto "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

func (t Endpoint) SendSecurityAlertEmail(ctx *gin.Context) {
	request := &corecontract.Request[emaildto.SendSecurityAlertEmailRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:  corecontract.Version,
			Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exceptions.New(
				"InvalidRequest",
				"Email",
				"SendSecurityAlertEmail",
				"The email request is invalid",
				http.StatusBadRequest,
			).WithOrigin(err),
		})
		return
	}
	if request.Operation != emaildto.SendSecurityAlertEmailOperation {
		ctx.JSON(http.StatusBadRequest, corecontract.Response[emaildto.SendEmailResponseDto]{
			Version:  corecontract.Version,
			Metadata: corecontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exceptions.New(
				"InvalidOperation",
				"Email",
				"SendSecurityAlertEmail",
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
				"SendSecurityAlertEmail",
				"The email request is invalid",
				http.StatusBadRequest,
			).WithOrigin(err),
		})
		return
	}

	exception := t.sender.SendSecurityAlertEmail(
		request.Dto.To,
		request.Dto.UserName,
		request.Dto.Status,
		request.Dto.AlertType,
		request.Dto.Reason,
		request.Dto.TimeOfOccurrence,
		request.Dto.OtherDetails,
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

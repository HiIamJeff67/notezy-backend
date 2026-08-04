package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	emaildto "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

func (t Endpoint) SendSecurityAlertEmail(ctx *gin.Context) {
	request := &gatewaycontract.Request[emaildto.SendSecurityAlertEmailRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gatewaycontract.Response[emaildto.SendEmailResponseDto]{
			Version:  gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
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
		ctx.JSON(http.StatusBadRequest, gatewaycontract.Response[emaildto.SendEmailResponseDto]{
			Version:  gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
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
	if err := t.validator.Struct(request.Dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gatewaycontract.Response[emaildto.SendEmailResponseDto]{
			Version:  gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
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
		ctx.JSON(exception.HTTPStatusCode(), gatewaycontract.Response[emaildto.SendEmailResponseDto]{
			Version:   gatewaycontract.Version,
			Metadata:  gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[emaildto.SendEmailResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     emaildto.SendEmailResponseDto{QueuedAt: time.Now()},
	})
}

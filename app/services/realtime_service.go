package services

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	tokens "github.com/HiIamJeff67/notezy-backend/app/tokens"
	validation "github.com/HiIamJeff67/notezy-backend/app/validation"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RealtimeServiceInterface interface {
	CreateMyRealtimeConnectionTicket(ctx context.Context, reqDto *dtos.CreateMyRealtimeConnectionTicketReqDto) (*dtos.CreateMyRealtimeConnectionTicketResDto, *exceptions.Exception)
	CreateMyBlockPackChannelTicket(ctx context.Context, reqDto *dtos.CreateMyBlockPackChannelTicketReqDto) (*dtos.CreateMyBlockPackChannelTicketResDto, *exceptions.Exception)
}

type RealtimeService struct {
	db                  *gorm.DB
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewRealtimeService(
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) RealtimeServiceInterface {
	return &RealtimeService{
		db:                  db,
		blockPackRepository: blockPackRepository,
	}
}

func (s *RealtimeService) CreateMyRealtimeConnectionTicket(
	ctx context.Context,
	reqDto *dtos.CreateMyRealtimeConnectionTicketReqDto,
) (*dtos.CreateMyRealtimeConnectionTicketResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Token.InvalidDto().WithOrigin(err)
	}

	connectionTicket, expiresAt, exception := tokens.GenerateRealtimeConnectionTicket(
		reqDto.ContextFields.UserPublicId,
		reqDto.Header.UserAgent,
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateMyRealtimeConnectionTicketResDto{
		RealtimeEndpoint:        "/" + constants.RealtimeDevelopmentBaseURL,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		ConnectionTicket:        *connectionTicket,
		ExpiresAt:               expiresAt,
	}, nil
}

func (s *RealtimeService) CreateMyBlockPackChannelTicket(
	ctx context.Context,
	reqDto *dtos.CreateMyBlockPackChannelTicketReqDto,
) (*dtos.CreateMyBlockPackChannelTicketResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	if reqDto.Body.Permission == realtimetypes.ChannelPermission_Read {
		allowedPermissions = append(
			allowedPermissions,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		)
	} else {
		allowedPermissions = append(allowedPermissions, enums.AccessControlPermission_Write)
	}

	blockPack, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		reqDto.Body.BlockPackId,
		reqDto.ContextFields.UserId,
		[]schemas.BlockPackRelation{schemas.BlockPackRelation_YjsDocument},
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	channelTicket, expiresAt, exception := tokens.GenerateRealtimeBlockPackTicket(
		reqDto.ContextFields.UserPublicId,
		reqDto.Header.UserAgent,
		blockPack.Id,
		reqDto.Body.Permission,
	)
	if exception != nil {
		return nil, exception
	}

	if blockPack.YjsDocument == nil {
		return nil, exceptions.BlockPack.NotFound("block pack yjs document is missing")
	}

	return &dtos.CreateMyBlockPackChannelTicketResDto{
		ChannelTicket:           *channelTicket,
		ExpiresAt:               expiresAt,
		ChannelType:             realtimetypes.ChannelType_BlockPack,
		ChannelId:               blockPack.Id,
		Permission:              reqDto.Body.Permission,
		RoomName:                fmt.Sprintf("%s:%s", constants.YjsBlockPackRoomPrefix, blockPack.Id),
		FragmentName:            constants.YjsBlockPackFragmentName,
		SchemaId:                constants.YjsBlockPackSchemaId,
		SchemaVersion:           constants.YjsBlockPackSchemaVersion,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		LastUpdateSequence:      blockPack.YjsDocument.LastUpdateSequence,
		CompactedUntilSequence:  blockPack.YjsDocument.CompactedUntilSequence,
	}, nil
}

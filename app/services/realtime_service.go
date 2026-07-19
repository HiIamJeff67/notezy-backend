package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
	GetBlockPackChannelAdmission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission realtimetypes.ChannelPermission) (int32, error)
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

func (s *RealtimeService) GetBlockPackChannelAdmission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission realtimetypes.ChannelPermission,
) (int32, error) {
	db := s.db.WithContext(ctx)

	var user schemas.User
	if err := db.
		Select("id").
		Where("public_id = ?", userPublicId).
		First(&user).Error; err != nil {
		return 0, err
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	if permission == realtimetypes.ChannelPermission_Read {
		allowedPermissions = append(
			allowedPermissions,
			enums.AccessControlPermission_Write,
			enums.AccessControlPermission_Read,
		)
	} else if permission == realtimetypes.ChannelPermission_Write {
		allowedPermissions = append(allowedPermissions, enums.AccessControlPermission_Write)
	} else {
		return 0, errors.New("invalid realtime channel permission")
	}

	ownerId, _, exception := s.blockPackRepository.CheckPermissionAndGetOneWithOwnerIdById(
		blockPackId,
		user.Id,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return 0, exception
	}

	var maximumSubscribers int32
	result := db.
		Model(&schemas.User{}).
		Select(`"PlanLimitationTable".max_realtime_room_subscriber_count`).
		Joins(`INNER JOIN "PlanLimitationTable" ON "PlanLimitationTable".key = "UserTable".plan`).
		Where(`"UserTable".id = ?`, *ownerId).
		Scan(&maximumSubscribers)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		return 0, errors.New("block pack owner has no realtime room subscriber capacity")
	}

	return maximumSubscribers, nil
}

package services

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"

	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
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
	GetMyBlockPackRealtimeParticipants(ctx context.Context, reqDto *dtos.GetMyBlockPackRealtimeParticipantsReqDto) (*dtos.GetMyBlockPackRealtimeParticipantsResDto, *exceptions.Exception)
	GetBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission realtimetypes.ChannelPermission) (int32, realtimetypes.ErrorCode, error)
	ValidateBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission realtimetypes.ChannelPermission) (realtimetypes.ErrorCode, error)
	CreateMyRealtimeConnectionTicket(ctx context.Context, reqDto *dtos.CreateMyRealtimeConnectionTicketReqDto) (*dtos.CreateMyRealtimeConnectionTicketResDto, *exceptions.Exception)
	CreateMyBlockPackChannelTicket(ctx context.Context, reqDto *dtos.CreateMyBlockPackChannelTicketReqDto) (*dtos.CreateMyBlockPackChannelTicketResDto, *exceptions.Exception)
}

type RealtimeService struct {
	db                  *gorm.DB
	blockPackRepository repositories.BlockPackRepositoryInterface
	leaseStore          *caches.RealtimeLeaseStore
}

func NewRealtimeService(
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) RealtimeServiceInterface {
	return &RealtimeService{
		db:                  db,
		blockPackRepository: blockPackRepository,
		leaseStore:          caches.NewRealtimeLeaseStore(caches.RedisClientMap),
	}
}

func (s *RealtimeService) GetMyBlockPackRealtimeParticipants(
	ctx context.Context, reqDto *dtos.GetMyBlockPackRealtimeParticipantsReqDto,
) (*dtos.GetMyBlockPackRealtimeParticipantsResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.BlockPack.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	_, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		reqDto.Param.BlockPackId,
		reqDto.ContextFields.UserId,
		nil,
		[]enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
		},
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	participants, err := s.leaseStore.GetBlockPackParticipants(reqDto.Param.BlockPackId)
	if err != nil {
		return nil, exceptions.Cache.FailedToUpdate("realtime block pack participants").WithOrigin(err)
	}
	if len(participants) == 0 {
		resDto := dtos.GetMyBlockPackRealtimeParticipantsResDto{}
		return &resDto, nil
	}

	connectionCountByPublicId := make(map[uuid.UUID]int, len(participants))
	permissionByPublicId := make(map[uuid.UUID]realtimetypes.ChannelPermission, len(participants))
	userPublicIds := make([]uuid.UUID, 0, len(participants))
	for _, participant := range participants {
		userPublicId, err := uuid.Parse(participant.UserPublicId)
		if err != nil {
			continue
		}

		if connectionCountByPublicId[userPublicId] == 0 {
			userPublicIds = append(userPublicIds, userPublicId)
		}
		connectionCountByPublicId[userPublicId]++
		permissionByPublicId[userPublicId] = realtimetypes.ChannelPermission(participant.ChannelPermission)
	}
	if len(userPublicIds) == 0 {
		resDto := dtos.GetMyBlockPackRealtimeParticipantsResDto{}
		return &resDto, nil
	}

	var users []schemas.User
	result := db.
		Model(&schemas.User{}).
		Select("public_id, name, display_name").
		Where("public_id IN ?", userPublicIds).
		Find(&users)
	if result.Error != nil {
		return nil, exceptions.User.NotFound().WithOrigin(result.Error)
	}

	resDto := make(dtos.GetMyBlockPackRealtimeParticipantsResDto, 0, len(users))
	for _, user := range users {
		resDto = append(resDto, dtos.RealtimeBlockPackParticipantResDto{
			UserPublicId:      user.PublicId,
			Name:              user.Name,
			DisplayName:       user.DisplayName,
			ChannelPermission: permissionByPublicId[user.PublicId],
			ConnectionCount:   connectionCountByPublicId[user.PublicId],
		})
	}

	sort.Slice(resDto, func(first int, second int) bool {
		return resDto[first].DisplayName < resDto[second].DisplayName
	})

	return &resDto, nil
}

func (s *RealtimeService) GetBlockPackChannelPermission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission realtimetypes.ChannelPermission,
) (int32, realtimetypes.ErrorCode, error) {
	db := s.db.WithContext(ctx)

	if errorCode, err := s.ValidateBlockPackChannelPermission(ctx, userPublicId, blockPackId, permission); err != nil {
		return 0, errorCode, err
	}

	var rootShelf schemas.RootShelf
	result := db.
		Model(&schemas.BlockPack{}).
		Select(`"RootShelfTable".owner_id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "RootShelfTable" ON "RootShelfTable".id = "SubShelfTable".root_shelf_id`).
		Where(`"BlockPackTable".id = ?`, blockPackId).
		Where(`"BlockPackTable".deleted_at IS NULL`).
		Where(`"SubShelfTable".deleted_at IS NULL`).
		Where(`"RootShelfTable".deleted_at IS NULL`).
		Scan(&rootShelf)
	if result.Error != nil {
		return 0, realtimetypes.ErrorCode_ResourceUnavailable, result.Error
	}
	if result.RowsAffected == 0 || rootShelf.OwnerId == uuid.Nil {
		return 0, realtimetypes.ErrorCode_ResourceUnavailable, gorm.ErrRecordNotFound
	}

	var maximumSubscribers int32
	result = db.
		Model(&schemas.User{}).
		Select(`"PlanLimitationTable".max_realtime_room_subscriber_count`).
		Joins(`INNER JOIN "PlanLimitationTable" ON "PlanLimitationTable".key = "UserTable".plan`).
		Where(`"UserTable".id = ?`, rootShelf.OwnerId).
		Scan(&maximumSubscribers)
	if result.Error != nil {
		return 0, realtimetypes.ErrorCode_RoomAdmissionUnavailable, result.Error
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		return 0, realtimetypes.ErrorCode_RoomAdmissionUnavailable, errors.New("block pack owner has no realtime room subscriber capacity")
	}

	return maximumSubscribers, "", nil
}

func (s *RealtimeService) ValidateBlockPackChannelPermission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission realtimetypes.ChannelPermission,
) (realtimetypes.ErrorCode, error) {
	db := s.db.WithContext(ctx)

	if permission != realtimetypes.ChannelPermission_Read &&
		permission != realtimetypes.ChannelPermission_Write {
		return realtimetypes.ErrorCode_PermissionRevoked, errors.New("invalid realtime channel permission")
	}

	var exists int
	result := db.
		Model(&schemas.BlockPack{}).
		Select(`1`).
		Joins(`INNER JOIN "BlockPackYjsDocumentTable" ON "BlockPackYjsDocumentTable".block_pack_id = "BlockPackTable".id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".id = "BlockPackTable".parent_sub_shelf_id`).
		Joins(`INNER JOIN "RootShelfTable" ON "RootShelfTable".id = "SubShelfTable".root_shelf_id`).
		Where(`"BlockPackTable".id = ?`, blockPackId).
		Where(`"BlockPackTable".deleted_at IS NULL`).
		Where(`"BlockPackYjsDocumentTable".deleted_at IS NULL`).
		Where(`"SubShelfTable".deleted_at IS NULL`).
		Where(`"RootShelfTable".deleted_at IS NULL`).
		Limit(1).
		Scan(&exists)
	if result.Error != nil {
		return realtimetypes.ErrorCode_ResourceUnavailable, result.Error
	}
	if result.RowsAffected == 0 {
		return realtimetypes.ErrorCode_ResourceUnavailable, gorm.ErrRecordNotFound
	}

	var user schemas.User
	if err := db.
		Select("id").
		Where("public_id = ?", userPublicId).
		First(&user).Error; err != nil {
		return realtimetypes.ErrorCode_PermissionRevoked, err
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
	}

	_, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		blockPackId,
		user.Id,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return realtimetypes.ErrorCode_PermissionRevoked, exception
	}

	return "", nil
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
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	var yjsDocument schemas.BlockPackYjsDocument
	result := db.
		Where("block_pack_id = ?", blockPack.Id).
		Where("deleted_at IS NULL").
		First(&yjsDocument)
	if result.Error != nil {
		return nil, exceptions.BlockPack.NotFound().WithOrigin(result.Error)
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
		LastUpdateSequence:      yjsDocument.LastUpdateSequence,
		CompactedUntilSequence:  yjsDocument.CompactedUntilSequence,
	}, nil
}

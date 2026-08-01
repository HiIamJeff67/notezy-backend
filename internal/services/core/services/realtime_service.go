package services

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/realtime"
	caches "github.com/HiIamJeff67/notezy-backend/internal/caches"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RealtimeServiceInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx context.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto) (*realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto, *exceptions.Exception)
	GetBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission sharedtypes.ChannelPermission) (int32, sharedtypes.ErrorCode, error)
	ValidateBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission sharedtypes.ChannelPermission) (sharedtypes.ErrorCode, error)
	CreateMyRealtimeConnectionTicket(ctx context.Context, requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto) (*realtimedto.CreateMyRealtimeConnectionTicketResponseDto, *exceptions.Exception)
	CreateMyBlockPackChannelTicket(ctx context.Context, requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto) (*realtimedto.CreateMyBlockPackChannelTicketResponseDto, *exceptions.Exception)
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

func (s *RealtimeService) getActorUserPublicId(ctx context.Context) (uuid.UUID, *exceptions.Exception) {
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return uuid.Nil, exception
	}
	var user schemas.User
	result := s.db.WithContext(ctx).
		Model(&schemas.User{}).
		Select("public_id").
		Where("id = ?", actorUserId).
		First(&user)
	if result.Error != nil {
		return uuid.Nil, exceptions.New(
			"NotFound",
			"User",
			"ResolveActor",
			"User was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}

	return user.PublicId, nil
}

func (s *RealtimeService) GetMyBlockPackRealtimeParticipants(
	ctx context.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
) (*realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"BlockPack",
			"GetMyBlockPackRealtimeParticipants",
			"Realtime participant request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	allowedPermissions, exception := contexts.GetAllowedPermissions(ctx)
	if exception != nil {
		return nil, exception
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}

	_, exception = s.blockPackRepository.CheckPermissionAndGetOneById(
		requestDto.Param.BlockPackId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return nil, exception
	}

	participants, err := s.leaseStore.GetBlockPackParticipants(requestDto.Param.BlockPackId)
	if err != nil {
		return nil, exceptions.New(
			"ReadFailed",
			"Cache",
			"GetMyBlockPackRealtimeParticipants",
			"Failed to retrieve realtime participants",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if len(participants) == 0 {
		responseDto := realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto{}
		return &responseDto, nil
	}

	connectionCountByPublicId := make(map[uuid.UUID]int, len(participants))
	permissionByPublicId := make(map[uuid.UUID]sharedtypes.ChannelPermission, len(participants))
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
		permissionByPublicId[userPublicId] = sharedtypes.ChannelPermission(participant.ChannelPermission)
	}
	if len(userPublicIds) == 0 {
		responseDto := realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto{}
		return &responseDto, nil
	}

	var users []schemas.User
	result := db.
		Model(&schemas.User{}).
		Select("public_id, name, display_name").
		Where("public_id IN ?", userPublicIds).
		Find(&users)
	if result.Error != nil {
		return nil, exceptions.New(
			"QueryFailed",
			"User",
			"GetMyBlockPackRealtimeParticipants",
			"Failed to retrieve realtime participants",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	responseDto := make(realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto, 0, len(users))
	for _, user := range users {
		responseDto = append(responseDto, realtimedto.RealtimeBlockPackParticipantResponseDto{
			UserPublicId:      user.PublicId,
			Name:              user.Name,
			DisplayName:       user.DisplayName,
			ChannelPermission: string(permissionByPublicId[user.PublicId]),
			ConnectionCount:   connectionCountByPublicId[user.PublicId],
		})
	}

	sort.Slice(responseDto, func(first int, second int) bool {
		return responseDto[first].DisplayName < responseDto[second].DisplayName
	})

	return &responseDto, nil
}

func (s *RealtimeService) GetBlockPackChannelPermission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission sharedtypes.ChannelPermission,
) (int32, sharedtypes.ErrorCode, error) {
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
		return 0, sharedtypes.ErrorCode_ResourceUnavailable, result.Error
	}
	if result.RowsAffected == 0 || rootShelf.OwnerId == uuid.Nil {
		return 0, sharedtypes.ErrorCode_ResourceUnavailable, gorm.ErrRecordNotFound
	}

	var maximumSubscribers int32
	result = db.
		Model(&schemas.User{}).
		Select(`"PlanLimitationTable".max_realtime_room_subscriber_count`).
		Joins(`INNER JOIN "PlanLimitationTable" ON "PlanLimitationTable".key = "UserTable".plan`).
		Where(`"UserTable".id = ?`, rootShelf.OwnerId).
		Scan(&maximumSubscribers)
	if result.Error != nil {
		return 0, sharedtypes.ErrorCode_RoomAdmissionUnavailable, result.Error
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		return 0, sharedtypes.ErrorCode_RoomAdmissionUnavailable, errors.New("block pack owner has no realtime room subscriber capacity")
	}

	return maximumSubscribers, "", nil
}

func (s *RealtimeService) ValidateBlockPackChannelPermission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission sharedtypes.ChannelPermission,
) (sharedtypes.ErrorCode, error) {
	if permission != sharedtypes.ChannelPermission_Read &&
		permission != sharedtypes.ChannelPermission_Write {
		return sharedtypes.ErrorCode_PermissionRevoked, errors.New("invalid realtime channel permission")
	}

	sharedAllowedPermissions := permission.AllowedAccessControlPermissions()
	if len(sharedAllowedPermissions) == 0 {
		return sharedtypes.ErrorCode_PermissionRevoked, errors.New("invalid realtime channel permission")
	}
	allowedPermissions := make([]enums.AccessControlPermission, len(sharedAllowedPermissions))
	for index, sharedAllowedPermission := range sharedAllowedPermissions {
		allowedPermissions[index] = enums.AccessControlPermission(sharedAllowedPermission)
	}

	ctx = contexts.WithAllowedPermissions(ctx, allowedPermissions)
	db := s.db.WithContext(ctx)

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
		return sharedtypes.ErrorCode_ResourceUnavailable, result.Error
	}
	if result.RowsAffected == 0 {
		return sharedtypes.ErrorCode_ResourceUnavailable, gorm.ErrRecordNotFound
	}

	var user schemas.User
	if err := db.
		Select("id").
		Where("public_id = ?", userPublicId).
		First(&user).Error; err != nil {
		return sharedtypes.ErrorCode_PermissionRevoked, err
	}

	_, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		blockPackId,
		user.Id,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return sharedtypes.ErrorCode_PermissionRevoked, exception
	}

	return "", nil
}

func (s *RealtimeService) CreateMyRealtimeConnectionTicket(
	ctx context.Context,
	requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto,
) (*realtimedto.CreateMyRealtimeConnectionTicketResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"Realtime",
			"CreateMyRealtimeConnectionTicket",
			"Realtime connection ticket request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	userPublicId, exception := s.getActorUserPublicId(ctx)
	if exception != nil {
		return nil, exception
	}
	userAgentHash := sha256.Sum256([]byte(requestDto.Header.UserAgent))
	connectionClaims := sharedtokens.RealtimeConnectionTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
	}
	connectionClaims.Subject = userPublicId.String()
	connectionTicket, expiresAt, err := sharedtokens.GenerateRealtimeConnectionTicket(connectionClaims)
	if err != nil {
		return nil, exceptions.New(
			"GenerationFailed",
			"Realtime",
			"CreateMyRealtimeConnectionTicket",
			"Failed to generate the realtime connection ticket",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &realtimedto.CreateMyRealtimeConnectionTicketResponseDto{
		RealtimeEndpoint:        "/" + constants.RealtimeDevelopmentBaseURL,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		ConnectionTicket:        *connectionTicket,
		ExpiresAt:               expiresAt,
	}, nil
}

func (s *RealtimeService) CreateMyBlockPackChannelTicket(
	ctx context.Context,
	requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto,
) (*realtimedto.CreateMyBlockPackChannelTicketResponseDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"BlockPack",
			"CreateMyBlockPackChannelTicket",
			"Block pack channel ticket request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	permission := sharedtypes.ChannelPermission(requestDto.Body.Permission)
	sharedAllowedPermissions := permission.AllowedAccessControlPermissions()
	if len(sharedAllowedPermissions) == 0 {
		return nil, exceptions.New(
			"InvalidChannelPermission",
			"BlockPack",
			"CreateMyBlockPackChannelTicket",
			"Realtime channel permission is invalid",
			http.StatusBadRequest,
		)
	}
	allowedPermissions := make([]enums.AccessControlPermission, len(sharedAllowedPermissions))
	for index, sharedAllowedPermission := range sharedAllowedPermissions {
		allowedPermissions[index] = enums.AccessControlPermission(sharedAllowedPermission)
	}
	actorUserId, exception := contexts.GetActorUserId(ctx)
	if exception != nil {
		return nil, exception
	}
	userPublicId, exception := s.getActorUserPublicId(ctx)
	if exception != nil {
		return nil, exception
	}

	blockPack, exception := s.blockPackRepository.CheckPermissionAndGetOneById(
		requestDto.Body.BlockPackId,
		actorUserId,
		nil,
		allowedPermissions,
		options.WithDB(db),
		options.WithAllowedPermissions(allowedPermissions),
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
		return nil, exceptions.New(
			"NotFound",
			"BlockPackYjsDocument",
			"CreateMyBlockPackChannelTicket",
			"Block pack document was not found",
			http.StatusNotFound,
		).WithOrigin(result.Error)
	}

	userAgentHash := sha256.Sum256([]byte(requestDto.Header.UserAgent))
	channelClaims := sharedtokens.RealtimeBlockPackTicketClaims{
		UserAgentHash:           fmt.Sprintf("%x", userAgentHash),
		ChannelType:             "BlockPack",
		ChannelId:               blockPack.Id.String(),
		Permission:              string(permission),
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		SchemaVersion:           constants.YjsBlockPackSchemaVersion,
	}
	channelClaims.Subject = userPublicId.String()
	channelTicket, expiresAt, err := sharedtokens.GenerateRealtimeBlockPackTicket(channelClaims)
	if err != nil {
		return nil, exceptions.New(
			"GenerationFailed",
			"BlockPack",
			"CreateMyBlockPackChannelTicket",
			"Failed to generate the block pack channel ticket",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &realtimedto.CreateMyBlockPackChannelTicketResponseDto{
		ChannelTicket:           *channelTicket,
		ExpiresAt:               expiresAt,
		ChannelType:             "BlockPack",
		ChannelId:               blockPack.Id,
		Permission:              string(permission),
		RoomName:                fmt.Sprintf("%s:%s", constants.YjsBlockPackRoomPrefix, blockPack.Id),
		FragmentName:            constants.YjsBlockPackFragmentName,
		SchemaId:                constants.YjsBlockPackSchemaId,
		SchemaVersion:           constants.YjsBlockPackSchemaVersion,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		LastUpdateSequence:      yjsDocument.LastUpdateSequence,
		CompactedUntilSequence:  yjsDocument.CompactedUntilSequence,
	}, nil
}

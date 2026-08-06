package realtime

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"sort"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"
	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"
	enumscontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	yjsworkercontract "github.com/HiIamJeff67/notezy-backend/contracts/yjsworker/v1"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	options "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas/enums"
)

type RealtimeServiceInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx context.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto) (*realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto, *exceptions.Exception)
	CreateMyRealtimeConnectionTicket(ctx context.Context, requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto) (*realtimedto.CreateMyRealtimeConnectionTicketResponseDto, *exceptions.Exception)
	CreateMyBlockPackChannelTicket(ctx context.Context, requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto) (*realtimedto.CreateMyBlockPackChannelTicketResponseDto, *exceptions.Exception)
}

type RealtimeService struct {
	validator           *validator.Validate
	db                  *gorm.DB
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewRealtimeService(
	validator *validator.Validate,
	db *gorm.DB,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) RealtimeServiceInterface {
	return &RealtimeService{
		validator:           validator,
		db:                  db,
		blockPackRepository: blockPackRepository,
	}
}

/* ============================== Auxiliary Functions ============================== */

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

func (s *RealtimeService) getBlockPackMaximumSubscribers(
	ctx context.Context,
	blockPackId uuid.UUID,
) (int32, error) {
	var rootShelf schemas.RootShelf
	result := s.db.WithContext(ctx).
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
		return 0, result.Error
	}
	if result.RowsAffected == 0 || rootShelf.OwnerId == uuid.Nil {
		return 0, gorm.ErrRecordNotFound
	}

	var maximumSubscribers int32
	result = s.db.WithContext(ctx).
		Model(&schemas.User{}).
		Select(`"PlanLimitationTable".max_realtime_room_subscriber_count`).
		Joins(`INNER JOIN "PlanLimitationTable" ON "PlanLimitationTable".key = "UserTable".plan`).
		Where(`"UserTable".id = ?`, rootShelf.OwnerId).
		Scan(&maximumSubscribers)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 || maximumSubscribers <= 0 {
		return 0, errors.New("block pack owner has no realtime room subscriber capacity")
	}

	return maximumSubscribers, nil
}

/* ============================== Service Methods for Realtime ============================== */

func (s *RealtimeService) GetMyBlockPackRealtimeParticipants(
	ctx context.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
) (*realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
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

	if len(requestDto.Body.Participants) == 0 {
		responseDto := realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto{}
		return &responseDto, nil
	}

	participantByPublicId := make(map[uuid.UUID]realtimedto.RealtimeBlockPackParticipantRequestDto, len(requestDto.Body.Participants))
	for _, participant := range requestDto.Body.Participants {
		existingParticipant, exists := participantByPublicId[participant.UserPublicId]
		if !exists {
			participantByPublicId[participant.UserPublicId] = participant
			continue
		}

		existingParticipant.ConnectionCount += participant.ConnectionCount
		if existingParticipant.ConnectionCount > constants.RealtimeMaxConnectionsPerUser {
			return nil, exceptions.New(
				"InvalidRequest",
				"BlockPack",
				"GetMyBlockPackRealtimeParticipants",
				"Realtime participant connection count is invalid",
				http.StatusBadRequest,
			)
		}
		if existingParticipant.ChannelPermission != "write" || participant.ChannelPermission == "write" {
			existingParticipant.ChannelPermission = participant.ChannelPermission
		}
		participantByPublicId[participant.UserPublicId] = existingParticipant
	}
	if len(participantByPublicId) > constants.MaxSearchLimit {
		return nil, exceptions.New(
			"InvalidRequest",
			"BlockPack",
			"GetMyBlockPackRealtimeParticipants",
			"Realtime participant request contains too many identities",
			http.StatusBadRequest,
		)
	}

	userPublicIds := make([]uuid.UUID, 0, len(participantByPublicId))
	for userPublicId := range participantByPublicId {
		userPublicIds = append(userPublicIds, userPublicId)
	}

	var users []schemas.User
	result := db.
		Model(&schemas.User{}).
		Select("public_id, name, display_name").
		Joins(`INNER JOIN "UsersToShelvesTable" AS users_to_shelves ON users_to_shelves.user_id = "UserTable".id`).
		Joins(`INNER JOIN "SubShelfTable" ON "SubShelfTable".root_shelf_id = users_to_shelves.root_shelf_id`).
		Joins(`INNER JOIN "BlockPackTable" ON "BlockPackTable".parent_sub_shelf_id = "SubShelfTable".id`).
		Where(`"BlockPackTable".id = ?`, requestDto.Param.BlockPackId).
		Where(`"BlockPackTable".deleted_at IS NULL`).
		Where(`"SubShelfTable".deleted_at IS NULL`).
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
		participant := participantByPublicId[user.PublicId]
		responseDto = append(responseDto, realtimedto.RealtimeBlockPackParticipantResponseDto{
			UserPublicId:      user.PublicId,
			Name:              user.Name,
			DisplayName:       user.DisplayName,
			ChannelPermission: participant.ChannelPermission,
			ConnectionCount:   participant.ConnectionCount,
		})
	}

	sort.Slice(responseDto, func(first int, second int) bool {
		return responseDto[first].DisplayName < responseDto[second].DisplayName
	})

	return &responseDto, nil
}

func (s *RealtimeService) CreateMyRealtimeConnectionTicket(
	ctx context.Context,
	requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto,
) (*realtimedto.CreateMyRealtimeConnectionTicketResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
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
		RealtimeEndpoint:        "/" + realtimegatewaycontract.RealtimeDevelopmentBaseURL,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		ConnectionTicket:        *connectionTicket,
		ExpiresAt:               expiresAt,
	}, nil
}

func (s *RealtimeService) CreateMyBlockPackChannelTicket(
	ctx context.Context,
	requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto,
) (*realtimedto.CreateMyBlockPackChannelTicketResponseDto, *exceptions.Exception) {
	if err := s.validator.Struct(requestDto); err != nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"BlockPack",
			"CreateMyBlockPackChannelTicket",
			"Block pack channel ticket request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	permission := enumscontract.ChannelPermission(requestDto.Body.Permission)
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

	maximumSubscribers, err := s.getBlockPackMaximumSubscribers(ctx, blockPack.Id)
	if err != nil {
		return nil, exceptions.New(
			"Unavailable",
			"BlockPack",
			"CreateMyBlockPackChannelTicket",
			"Block pack realtime room admission is unavailable",
			http.StatusServiceUnavailable,
		).WithOrigin(err)
	}

	userAgentHash := sha256.Sum256([]byte(requestDto.Header.UserAgent))
	channelClaims := sharedtokens.RealtimeBlockPackTicketClaims{
		UserAgentHash:                    fmt.Sprintf("%x", userAgentHash),
		ChannelType:                      "BlockPack",
		ChannelId:                        blockPack.Id.String(),
		Permission:                       string(permission),
		RealtimeProtocolVersion:          constants.RealtimeProtocolVersion,
		SchemaVersion:                    yjsworkercontract.YjsBlockPackSchemaVersion,
		RoomAdmissionPolicyVersion:       realtimegatewaycontract.BlockPackRoomAdmissionPolicyVersion,
		RoomAdmissionEnforcementStrategy: string(realtimegatewaycontract.RoomAdmissionEnforcementStrategy_RejectNewSubscriber),
		MaximumSubscribers:               maximumSubscribers,
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
		RoomName:                fmt.Sprintf("%s:%s", yjsworkercontract.YjsBlockPackRoomPrefix, blockPack.Id),
		FragmentName:            yjsworkercontract.YjsBlockPackFragmentName,
		SchemaId:                yjsworkercontract.YjsBlockPackSchemaId,
		SchemaVersion:           yjsworkercontract.YjsBlockPackSchemaVersion,
		RealtimeProtocolVersion: constants.RealtimeProtocolVersion,
		LastUpdateSequence:      yjsDocument.LastUpdateSequence,
		CompactedUntilSequence:  yjsDocument.CompactedUntilSequence,
	}, nil
}

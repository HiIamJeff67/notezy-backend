package caches

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/google/uuid"

	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func TestRealtimeLeaseStoreLimitsConcurrentUserConnections(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start test redis server: %v", err)
	}
	defer server.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer redisClient.Close()

	store := NewRealtimeLeaseStore(map[int]*redis.Client{constants.RealtimeRedisServerNumber: redisClient})
	userPublicId := uuid.New()

	var acquiredCount atomic.Int32
	var waitGroup sync.WaitGroup
	for range 32 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			acquired, _, err := store.AcquireUserConnection(userPublicId, uuid.New(), 1)
			if err != nil {
				t.Errorf("failed to acquire realtime user connection lease: %v", err)

				return
			}
			if acquired {
				acquiredCount.Add(1)
			}
		}()
	}
	waitGroup.Wait()

	if acquiredCount.Load() != 1 {
		t.Fatalf("expected exactly one concurrent lease acquisition, got %d", acquiredCount.Load())
	}
}

func TestRealtimeLeaseStoreReclaimsExpiredUserConnectionLease(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start test redis server: %v", err)
	}
	defer server.Close()

	firstRedisClient := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer firstRedisClient.Close()
	secondRedisClient := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer secondRedisClient.Close()

	firstStore := NewRealtimeLeaseStore(map[int]*redis.Client{constants.RealtimeRedisServerNumber: firstRedisClient})
	secondStore := NewRealtimeLeaseStore(map[int]*redis.Client{constants.RealtimeRedisServerNumber: secondRedisClient})
	userPublicId := uuid.New()

	acquired, _, err := firstStore.AcquireUserConnection(userPublicId, uuid.New(), 1)
	if err != nil || !acquired {
		t.Fatalf("expected first user connection lease to be acquired: %v", err)
	}

	acquired, _, err = secondStore.AcquireUserConnection(userPublicId, uuid.New(), 1)
	if err != nil {
		t.Fatalf("failed to check second user connection lease: %v", err)
	}
	if acquired {
		t.Fatal("expected the second Redis client to observe the distributed user connection cap")
	}

	server.FastForward(constants.RealtimeLeaseTTL)

	acquired, _, err = secondStore.AcquireUserConnection(userPublicId, uuid.New(), 1)
	if err != nil || !acquired {
		t.Fatalf("expected expired user connection lease to be reclaimed: %v", err)
	}
}

func TestRealtimeLeaseStoreReleasesBlockPackSubscriber(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start test redis server: %v", err)
	}
	defer server.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer redisClient.Close()

	store := NewRealtimeLeaseStore(map[int]*redis.Client{constants.RealtimeRedisServerNumber: redisClient})
	blockPackId := uuid.New()

	acquired, _, err := store.AcquireBlockPackSubscriber(blockPackId, "connector-a:1", 1)
	if err != nil || !acquired {
		t.Fatalf("expected first BlockPack subscriber lease to be acquired: %v", err)
	}

	acquired, _, err = store.AcquireBlockPackSubscriber(blockPackId, "connector-b:1", 1)
	if err != nil {
		t.Fatalf("failed to check BlockPack subscriber lease: %v", err)
	}
	if acquired {
		t.Fatal("expected BlockPack subscriber capacity to reject the second lease")
	}

	if err := store.ReleaseBlockPackSubscriber(blockPackId, "connector-a:1"); err != nil {
		t.Fatalf("failed to release BlockPack subscriber lease: %v", err)
	}

	acquired, _, err = store.AcquireBlockPackSubscriber(blockPackId, "connector-b:1", 1)
	if err != nil || !acquired {
		t.Fatalf("expected released BlockPack subscriber capacity to be reusable: %v", err)
	}
}

func TestRealtimeLeaseStoreListsOnlyActiveBlockPackParticipants(t *testing.T) {
	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start test redis server: %v", err)
	}
	defer server.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: server.Addr()})
	defer redisClient.Close()

	store := NewRealtimeLeaseStore(map[int]*redis.Client{constants.RealtimeRedisServerNumber: redisClient})
	blockPackId := uuid.New()
	userPublicId := uuid.New()
	member := "connector-a:1"

	acquired, _, err := store.AcquireBlockPackSubscriber(blockPackId, member, 1)
	if err != nil || !acquired {
		t.Fatalf("expected BlockPack subscriber lease to be acquired: %v", err)
	}
	if err := store.SetBlockPackParticipant(blockPackId, member, userPublicId, "write"); err != nil {
		t.Fatalf("failed to record BlockPack participant: %v", err)
	}

	participants, err := store.GetBlockPackParticipants(blockPackId)
	if err != nil {
		t.Fatalf("failed to list BlockPack participants: %v", err)
	}
	if len(participants) != 1 || participants[0].UserPublicId != userPublicId.String() ||
		participants[0].ChannelPermission != "write" {
		t.Fatalf("unexpected BlockPack participants: %#v", participants)
	}

	if err := store.ReleaseBlockPackSubscriber(blockPackId, member); err != nil {
		t.Fatalf("failed to release BlockPack subscriber lease: %v", err)
	}

	participants, err = store.GetBlockPackParticipants(blockPackId)
	if err != nil {
		t.Fatalf("failed to list BlockPack participants after release: %v", err)
	}
	if len(participants) != 0 {
		t.Fatalf("expected no active BlockPack participants after release, got %#v", participants)
	}
}

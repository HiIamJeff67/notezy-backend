package services

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	accountingtriggersql "github.com/HiIamJeff67/notezy-backend/app/models/schemas/triggers/accounting_triggers"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func TestBlockProjectionServiceAppliesCurrentDocumentState(t *testing.T) {
	databaseDSN := os.Getenv("NOTEZY_BLOCK_PROJECTION_INTEGRATION_DATABASE_DSN")
	if databaseDSN == "" {
		t.Skip("NOTEZY_BLOCK_PROJECTION_INTEGRATION_DATABASE_DSN is not configured; it must point to an isolated database because this test recreates projection tables")
	}

	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to block projection integration database: %v", err)
	}

	for _, statement := range []string{
		`DROP TABLE IF EXISTS "BlockTable"`,
		`DROP TABLE IF EXISTS "BlockPackYjsDocumentTable"`,
		`DROP TABLE IF EXISTS "UsersToShelvesTable"`,
		`DROP TABLE IF EXISTS "PlanLimitationTable"`,
		`DROP TABLE IF EXISTS "UserAccountTable"`,
		`DROP TABLE IF EXISTS "UserTable"`,
		`DROP TABLE IF EXISTS "SubShelfTable"`,
		`DROP TABLE IF EXISTS "BlockPackTable"`,
		`DROP TYPE IF EXISTS "AccessControlPermission"`,
		`DROP TYPE IF EXISTS "UserPlan"`,
		`DROP TYPE IF EXISTS "BlockType"`,
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,
		`CREATE TYPE "AccessControlPermission" AS ENUM ('Owner')`,
		`CREATE TYPE "UserPlan" AS ENUM ('Free')`,
		`CREATE TYPE "BlockType" AS ENUM ('paragraph')`,
		`CREATE TABLE "BlockPackTable" (
			id uuid PRIMARY KEY,
			parent_sub_shelf_id uuid NOT NULL,
			block_count bigint NOT NULL DEFAULT 0,
			deleted_at timestamptz NULL,
			updated_at timestamptz NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE "BlockPackYjsDocumentTable" (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			block_pack_id uuid NOT NULL UNIQUE,
			snapshot bytea NOT NULL,
			state_vector bytea NOT NULL,
			last_update_sequence bigint NOT NULL DEFAULT 0,
			compacted_until_sequence bigint NOT NULL DEFAULT 0,
			projected_until_sequence bigint NOT NULL DEFAULT -1,
			deleted_at timestamptz NULL,
			updated_at timestamptz NOT NULL DEFAULT now(),
			created_at timestamptz NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE "BlockTable" (
			id uuid PRIMARY KEY,
			block_pack_id uuid NOT NULL,
			parent_block_id uuid NULL,
			prev_block_id uuid NULL,
			next_block_id uuid NULL,
			type "BlockType" NOT NULL,
			props jsonb NOT NULL,
			content jsonb NULL,
			updated_at timestamptz NOT NULL DEFAULT now(),
			created_at timestamptz NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE "SubShelfTable" (
			id uuid PRIMARY KEY,
			root_shelf_id uuid NOT NULL
		)`,
		`CREATE TABLE "UserTable" (
			id uuid PRIMARY KEY,
			plan "UserPlan" NOT NULL
		)`,
		`CREATE TABLE "UserAccountTable" (
			user_id uuid PRIMARY KEY,
			block_count bigint NOT NULL,
			updated_at timestamptz NOT NULL DEFAULT now()
		)`,
		`CREATE TABLE "PlanLimitationTable" (
			key "UserPlan" PRIMARY KEY,
			max_block_count integer NOT NULL,
			max_block_count_per_block_pack integer NOT NULL
		)`,
		`CREATE TABLE "UsersToShelvesTable" (
			user_id uuid NOT NULL,
			root_shelf_id uuid NOT NULL,
			permission "AccessControlPermission" NOT NULL
		)`,
	} {
		if err := db.Exec(statement).Error; err != nil {
			t.Fatalf("failed to set up block projection integration database: %v", err)
		}
	}
	for _, triggerSQL := range []string{
		accountingtriggersql.AccountingInsertedBlockTriggerSQL,
		accountingtriggersql.AccountingDeletedBlockTriggerSQL,
	} {
		for _, statement := range strings.Split(triggerSQL, constants.SQLSeparator) {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if err := db.Exec(statement).Error; err != nil {
				t.Fatalf("failed to set up block projection accounting trigger: %v", err)
			}
		}
	}

	blockPackId := uuid.New()
	rootShelfId := uuid.New()
	parentSubShelfId := uuid.New()
	userId := uuid.New()
	removedBlockId := uuid.New()
	restoredBlockId := uuid.New()
	newBlockId := uuid.New()

	for _, statement := range []struct {
		SQL  string
		Args []any
	}{
		{`INSERT INTO "BlockPackTable" (id, parent_sub_shelf_id) VALUES (?, ?)`, []any{blockPackId, parentSubShelfId}},
		{`INSERT INTO "BlockPackYjsDocumentTable" (block_pack_id, snapshot, state_vector, last_update_sequence) VALUES (?, '', '', 3)`, []any{blockPackId}},
		{`INSERT INTO "SubShelfTable" (id, root_shelf_id) VALUES (?, ?)`, []any{parentSubShelfId, rootShelfId}},
		{`INSERT INTO "UserTable" (id, plan) VALUES (?, 'Free')`, []any{userId}},
		{`INSERT INTO "UserAccountTable" (user_id, block_count) VALUES (?, 10)`, []any{userId}},
		{`INSERT INTO "PlanLimitationTable" (key, max_block_count, max_block_count_per_block_pack) VALUES ('Free', 100, 100)`, nil},
		{`INSERT INTO "UsersToShelvesTable" (user_id, root_shelf_id, permission) VALUES (?, ?, 'Owner')`, []any{userId, rootShelfId}},
		{`INSERT INTO "BlockTable" (id, block_pack_id, type, props, content) VALUES (?, ?, 'paragraph', '{}', '[]')`, []any{removedBlockId, blockPackId}},
		{`INSERT INTO "BlockTable" (id, block_pack_id, type, props, content) VALUES (?, ?, 'paragraph', '{}', '[]')`, []any{restoredBlockId, blockPackId}},
	} {
		if err := db.Exec(statement.SQL, statement.Args...).Error; err != nil {
			t.Fatalf("failed to seed block projection integration database: %v", err)
		}
	}

	payload := []byte(`[
		{
			"id": "` + restoredBlockId.String() + `",
			"type": "paragraph",
			"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
			"content": [],
			"children": [
				{
					"id": "` + newBlockId.String() + `",
					"type": "paragraph",
					"props": {"backgroundColor":"default","textColor":"default","textAlignment":"left"},
					"content": [],
					"children": []
				}
			]
		}
	]`)

	var roots []dtos.ArborizedEditableBlock
	if err := json.Unmarshal(payload, &roots); err != nil {
		t.Fatalf("failed to unmarshal projected block forest: %v", err)
	}

	service := NewBlockProjectionService(db, adapters.NewEditableBlockAdapter())
	result, err := service.Apply(context.Background(), blockPackId, dtos.ApplyBlockProjectionInput{
		SchemaId:          "notezy.blocknote",
		SchemaVersion:     1,
		ProjectedSequence: 3,
		Blocks:            roots,
	})
	if err != nil {
		t.Fatalf("failed to apply block projection: %v", err)
	}
	if !result.Applied || result.ProjectedUntilSequence != 3 {
		t.Fatalf("unexpected projection result: %#v", result)
	}

	var blockPackCount int64
	if err := db.Raw(`SELECT block_count FROM "BlockPackTable" WHERE id = ?`, blockPackId).Scan(&blockPackCount).Error; err != nil {
		t.Fatalf("failed to read projected block pack count: %v", err)
	}
	if blockPackCount != 2 {
		t.Fatalf("expected block pack count 2, got %d", blockPackCount)
	}

	var userAccountCount int64
	if err := db.Raw(`SELECT block_count FROM "UserAccountTable" WHERE user_id = ?`, userId).Scan(&userAccountCount).Error; err != nil {
		t.Fatalf("failed to read projected user account count: %v", err)
	}
	if userAccountCount != 12 {
		t.Fatalf("expected user account count 12, got %d", userAccountCount)
	}

	var documentSequence int64
	if err := db.Raw(`SELECT projected_until_sequence FROM "BlockPackYjsDocumentTable" WHERE block_pack_id = ?`, blockPackId).Scan(&documentSequence).Error; err != nil {
		t.Fatalf("failed to read projection checkpoint: %v", err)
	}
	if documentSequence != 3 {
		t.Fatalf("expected projected sequence 3, got %d", documentSequence)
	}

	var blocks []struct {
		Id            uuid.UUID  `gorm:"column:id"`
		ParentBlockId *uuid.UUID `gorm:"column:parent_block_id"`
	}
	if err := db.Raw(`SELECT id, parent_block_id FROM "BlockTable" ORDER BY id`).Scan(&blocks).Error; err != nil {
		t.Fatalf("failed to read projected blocks: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 physical blocks, got %d", len(blocks))
	}

	blocksById := make(map[uuid.UUID]struct {
		ParentBlockId *uuid.UUID
	}, len(blocks))
	for _, block := range blocks {
		blocksById[block.Id] = struct {
			ParentBlockId *uuid.UUID
		}{
			ParentBlockId: block.ParentBlockId,
		}
	}
	if _, exists := blocksById[removedBlockId]; exists {
		t.Fatalf("expected removed block to be physically deleted")
	}
	if blocksById[newBlockId].ParentBlockId == nil || *blocksById[newBlockId].ParentBlockId != restoredBlockId {
		t.Fatalf("expected new block parent to be restored root")
	}

	staleResult, err := service.Apply(context.Background(), blockPackId, dtos.ApplyBlockProjectionInput{
		SchemaId:          "notezy.blocknote",
		SchemaVersion:     1,
		ProjectedSequence: 3,
		Blocks:            roots,
	})
	if err != nil {
		t.Fatalf("failed to retry current block projection: %v", err)
	}
	if staleResult.Applied || staleResult.ProjectedUntilSequence != 3 {
		t.Fatalf("expected idempotent stale projection result, got %#v", staleResult)
	}
}

-- 0000_plan_limitation_seed.up.sql
INSERT INTO "PlanLimitationTable" (
    key, 
    max_root_shelf_count, 
    max_block_pack_count, 
    max_block_count, 
    max_sync_block_count,
    max_material_count, 
    max_work_flow_count, 
    max_additional_item_count,
    max_sub_shelf_count_per_root_shelf,
    max_item_count_per_root_shelf,
    max_block_count_per_block_pack,
    max_material_size,
    updated_at,
    created_at
) VALUES
('Free',        10,     20,     1000,   10,     10,     2,      5,      20,     20,     100,    5242880,      NOW(), NOW()),
('Pro',         50,     100,    5000,   50,     50,     10,     50,     100,    100,    200,    20971520,     NOW(), NOW()),
('Premium',     150,    300,    15000,  150,    150,    30,     150,    200,    200,    500,    52428800,     NOW(), NOW()),
('Ultimate',    300,    200,    30000,  300,    300,    60,     300,    500,    500,    1000,   209715200,    NOW(), NOW()),
('Enterprise',  1000,   2000,   100000, 1000,   1000,   100,    1000,   1000,   1000,   1000,   524288000,    NOW(), NOW())
ON CONFLICT (key) DO UPDATE SET
    max_root_shelf_count = EXCLUDED.max_root_shelf_count, 
    max_block_pack_count = EXCLUDED.max_block_pack_count, 
    max_block_count = EXCLUDED.max_block_count, 
    max_sync_block_count = EXCLUDED.max_sync_block_count,
    max_material_count = EXCLUDED.max_material_count, 
    max_work_flow_count = EXCLUDED.max_work_flow_count, 
    max_additional_item_count = EXCLUDED.max_additional_item_count,
    max_sub_shelf_count_per_root_shelf = EXCLUDED.max_sub_shelf_count_per_root_shelf,
    max_item_count_per_root_shelf = EXCLUDED.max_item_count_per_root_shelf,
    max_block_count_per_block_pack = EXCLUDED.max_block_count_per_block_pack,
    max_material_size = EXCLUDED.max_material_size,
    updated_at = NOW();
# Backend Name Uniqueness Changelog

## 2026-07-02

Backend no longer enforces unique display names for these records:

- `RootShelf`: `owner_id + name`
- `SubShelf`: `root_shelf_id + path + name`
- `BlockPack`: `parent_sub_shelf_id + name`
- `Material`: `parent_sub_shelf_id + name`
- `Station`: `name`

The old GORM unique constraints are preserved as inline comments next to the affected fields.

Existing databases that already created these unique indexes must be rebuilt or have those indexes dropped manually. Removing the GORM tags only prevents future schema creation from adding them.

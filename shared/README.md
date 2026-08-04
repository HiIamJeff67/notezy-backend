# Shared

Cross-runtime helpers live here, beside `internal/` and `contracts/`.

- `lib/` is portable and must not import Notezy code.
- Other packages may provide shared application utilities when more than one runtime needs them.
- Domain business logic remains with its owning service.

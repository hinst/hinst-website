I've reviewed the task file and the relevant code. Here's what I found and my proposed plan.

## Current state

**Backend** (Go, stdlib `net/http`):
- 12 API endpoints, all registered via `namedWebFunction` slices in `webAppGoals.init()` and `webAppAdmin.init()` (see `server/webAppGoals.go`, `server/webAppAdmin.go`)
- 4 response models in `server/rest_objects`: `GoalObject`, `GoalPostObject`, `GoalPostHeader`, `GoalPostSearchIndexingHeader` (with an embedded-struct `tstype:",extends"` hack for tygo)
- Errors are `webError{Message, Status}` returned as JSON via the panic/recover wrapper in `server/webApp.go`
- Binary image endpoints: `/goal/image`, `/goalPost/image`
- Two POST endpoints reading raw text body: `/goalPost/setText`, `/goalPost/setTitle`

**Frontend**:
- `tygo.yaml` generates `frontend/src/typescript/generated/rest_objects.ts` (types only)
- 8 files import from it; hand-written `*Ex` wrappers (`goalObjectEx.ts`, `goalPostHeaderEx.ts`) and `apiClient.ts` sit on top
- `generate-types.bat` installs + runs tygo

## Proposed plan

### Phase 1 — Generate OpenAPI from Go with swaggo
1. Add `github.com/swaggo/swag` as a tool dep; add `//go:generate swag init -g server/webAppGoals.go --output docs ...` in `backend/main.go` (replacing the tygo directive). All handlers live in the `server` package, so one init covers everything.
2. Add swag annotations to the handlers in `webAppGoals.go` / `webAppAdmin.go`:
   - method + path (e.g. `// @Router /hinst-website/api/goalPosts [get]` — path prefix baked in via the `--server`/general info)
   - query params (`id`, `goalId`, `postDateTime`, `lang`, `query`, `index`, `isPublic`, `enabled`, `languageTag`) with types/descriptions
   - responses: reuse `rest_objects.*` types (swag reflects on the existing json tags), plus a shared `webError` component for 400/401/403/404/500
   - binary endpoints documented as `application/octet-stream` (real content-type is dynamic: `image/png` etc.)
   - `setText`/`setTitle` documented as `text/plain` request body
3. Produce `backend/docs/openapi.json` (+ optional `openapi.md` via `swag init --markdown` or Redoc HTML for browser viewing).
4. Remove the `tstype:",extends"` comment from `rest_objects` (swag flattens/embeds the embedded struct; TS output stays compatible).

### Phase 2 — Replace tygo with OpenAPI → TypeScript
Recommended: **`openapi-typescript`** (types-only, zero runtime deps, minimal diff).
1. Add `openapi-typescript` as a frontend devDependency.
2. Generate into the **same file path** `frontend/src/typescript/generated/rest_objects.ts` so all 8 import sites and the `*Ex` wrappers keep compiling unchanged — this is a drop-in replacement for the tygo output.
3. Rewrite `generate-types.bat`: install `swag` if missing → `swag init` in `backend` → `npx openapi-typescript` in `frontend`.
4. Delete `tygo.yaml`, the `//go:generate tygo generate` directive, and any tygo install references (README etc.).

Alternative (more invasive): **orval** to also generate a typed API client, replacing the hand-written `apiClient.ts`. I'd skip this for now — it rewrites ~11 consumer files and isn't required by the task ("replace Tygo" only mandates the types).

### Phase 3 — Verification
- Run `generate-types.bat`; diff new `rest_objects.ts` against the tygo-generated one — should be nearly identical (same json tags), confirming drop-in compatibility
- `go build ./...` in backend; validate spec covers all 12 endpoints
- `npm run build` (tsc) in frontend to confirm all imports still type-check
- Sanity-check the spec renders (Redoc) for humans

### Open questions
1. **Spec location**: `backend/docs/openapi.json` (swag default) — OK, or prefer committed `openapi.yaml`?
	1. Answer: Use YAML format. Save into file backend/OpenAPI.yaml
2. **TS generator**: go with `openapi-typescript` (recommended, minimal), or do you want orval/full client generation too?
	1. Answer: use `openapi-typescript`
3. The frontend `apiClient.ts` also references endpoints that **don't exist in this backend** (`/riddles/new`, `/riddles/primeNumbers`, `/pingUrlManually` — dead code, likely from a template). Leave them as-is, or clean up while we're at it?
	1. Answer: riddles are already deleted.
	1. pingUrlManually: delete it and replace old calls with //TODO

# Use OpenAPI type definitions for Frontend

* Backend generates OpenAPI spec, saves it to file backend/schema.yaml
	* Already done
* Frontend should generate TypeScript from OpenAPI spec
	* Already done
* Use package openapi-fetch in frontend
	* Done
	* `frontend/src/typescript/generated/rest_objects.ts` (leftover from the old Tygo approach) has been deleted
	* `frontend/src/typescript/apiClient.ts` now uses openapi-fetch with the generated types (`generated/openapi.d.ts`)
	* Type aliases for the old object types live in `frontend/src/typescript/apiTypes.ts`

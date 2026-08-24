# Use OpenAPI type definitions for Frontend

* Backend generates OpenAPI spec, saves it to file backend/schema.yaml
	* Already done
* Frontend should generate TypeScript from OpenAPI spec
	* Already done
* Use package openapi-fetch in frontend
	* This is not done yet
	* Right now frontend uses
		* `frontend/src/typescript/apiClient.ts`
		* `frontend/src/typescript/generated/rest_objects.ts`
			* This is actually no longer generated, it is leftover from the old approach using Tygo
	* The task is to delete old rest_object.ts and use openapi-fetch from now on

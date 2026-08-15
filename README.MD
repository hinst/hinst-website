# hinst-website

See `README.md` in sub-folders `frontend` and `backend` for more information.
See `CHANGELOG.md` to see what is new.

## License
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/

## How to deploy
1. Call `build && push` in the root folder to build Docker image and push it to registry
	1. Call `deploy` in `orange-pi-kubernetes` folder to deploy backend container to my Orange-Pi single board computer
2. Call `node build.mjs` in `frontend` folder
	1. Commit and push `hinst.github.io` repository to make frontend available on GitHub pages

## API type generation
The REST API object types shared between backend and frontend are defined in Golang in `backend/server/db_objects`. TypeScript types are generated from them with [tygo](https://github.com/gzuidhof/tygo):
1. Call `generate-types.bat` in the root folder (requires Go, tygo is installed automatically if missing)
2. Generated file: `frontend/src/typescript/generated/db_objects.ts` (do not edit manually)
3. Config: `backend/tygo.yaml`

## Administrator mode
In administrator mode we can edit article texts and toggle article visibility: private or public. Initially all newly downloaded articles are private.
1. To use administrator mode, set cookies:
	1. `goalManagerMode=1`
	2. `adminPassword=` actual password
2. Administrator mode is available only when API hostname is the same as the web page hostname. There is no easy way to use administrator mode in public proxied mode, because the domain name for the cookies will be different.

## Static mode
1. Call `backend\batch\generate.bat` to create static HTML files in `backend\static` folder
2. To deploy static files, copy them to hinst.github.io, commit and push


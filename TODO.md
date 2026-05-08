1. Move Smart Progress importer into backend using Go programming language
	1. Afterwards, deprecate Smart Progress importer: keep as archived repository, no longer used
	1. This change should reduce maintenance effort because maintaining one repository should be easier than maintaining two separate programs
1. Import additional information from Smart Progress
	1. Cover image
	1. Description
1. Use new model for AI translations: Gemma-4 instead of Aya-expanse
1. Use CPU only mode for model, to avoid video memory usage in background
1. Use markdown for storing blog posts
1. `long-term` Reverse blog import direction
	1. Before: using Smart Progress as primary source and hinst-website as secondary storage
	1. After: use hinst-website as primary source and Smart Progress for publishing
1. Better structure for source code files: group by package
1. Improvements for Kubernetes
	1. ~~Upgrade Kubernetes to latest version 1.36 on Orange Pi~~
	1. Put all required deployments into the same namespace `hinst-website`
		1. PostgreSQL
		1. Prettier-Server
	1. Set up and use my own container registry, instead of Docker.io
1. Automate dynamic website deployment
1. Show comments
1. Allow creating new comments
1. Put translated texts into separate tables

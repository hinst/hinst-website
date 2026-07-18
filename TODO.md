1. Move Smart Progress importer into backend using Go programming language
	1. Afterwards, deprecate Smart Progress importer: keep as archived repository, no longer used
	1. This change should reduce maintenance effort because maintaining one repository should be easier than maintaining two separate programs
1. Introduce reusable structure for storing translated texts, instead of copy-pasting `title`, `titleEnglish`, `titleGerman`
1. `long-term` Reverse blog import direction
	1. Before: using Smart Progress as primary source and hinst-website as secondary storage
	1. After: use hinst-website as primary source and Smart Progress for publishing
1. Automate dynamic website deployment
1. Show comments
1. Allow creating new comments
1. Put translated texts into separate tables
1. Run performance tests to see how many requests per second my website can handle

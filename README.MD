# Changelog

## 2025.07

* Moved source files to separated folders:
	* src/typescript
	* src/tsx
	* src/css
* Refactoring: merged LanguageContext and DisplayWidthContext into one AppContext
* Storage migration: changed from JSON files and folders to SQLite database
* Enabled prettier formatting for source code files
* Fixed bug after storage migration: provide isPublic
* Fixed back browser button: go to main page and skip the intermediate loading page
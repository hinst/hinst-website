# hinst-website

## How to deploy
1. Call `build && push` in the root folder
	1. Call `deploy` in `orange-pi-kubernetes` folder to deploy backend container to my Orange-Pi single board computer
2. Call `node build.js` in `frontend` folder
	1. Commit and push `hinst.github.io` repository to make frontend available on GitHub pages

## Changelog

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
* Fixed article layout: fill width but no overflow
* Added administrator feature: edit article text
	* The text editor is supposed to be used for correcting minor mistakes made by the LLM translator, such as: translating "e-scooter" as "skateboard"
* Added settings page
	* Color theme selection: System, Dark, Light
	* Language selection: System, English, Russian, German
* Fixed setInterval handling: avoid using captured state variables in inner functions
* Fixed timezone handling in downloader
* Cut href redirects from links in article texts
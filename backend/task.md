# Task

## Current status

Right now we generate static HTML website by launching web server and sending HTTP requests to it. Afterward, we have response into HTML files.

## Goal

Generate static HTML website directly from the database, bypassing HTTP protocol. Web server for static HTML pages is not needed anymore.

## Relevant files

* server/webPageGoals.go
* server/webStaticGoals.go

## Approach

Split work in small tasks and delegate them to the relevant subagent. Do NOT attempt to do all at once, or you will get brain explode.

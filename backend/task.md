# Task

Looking at files:
* goalPostRow.go
* schema.postgre.sql

See that these fields are nullable:
* TextEnglish
* TextGerman
* Title
* TitleEnglish
* TitleGerman

The goal of the task is to simplify our code and the database schema. These fields should become no longer nullable.
Instead of null, use empty string.

For translated text generator: instead of checking for null, check if target field is empty and the source field is non-empty.

You can either migrate all fields in one go, or launch a worker subagent for every field. It is up to you.

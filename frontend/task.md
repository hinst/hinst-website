# Use generated REST API response definitions

Earlier object definitions for REST API responses were hand-written, copying their `struct` counterparts in Go backend.
Recently we have set up code generator, saving REST API objects to `src\typescript\generated\rest_objects.ts`.
However new definitions are not used yet. We still use old definitions from `src\typescript\personal-goals`.
The goal is to use generated definitions and delete old definitions.
One obstacle is that old objects have methods such as `yearAndMonthText()`.
Please find a solution how to attach those methods to the new generated objects, but editing generated code is not allowed, so you have to use extends/implements or something.

# Rest object update

See files:
* src\typescript\restObjectExtensions.ts
* src\typescript\generated\rest_objects.ts
	* Generated file, editing not possible.

Change the approach used for REST object extension methods.
Instead of `declare module` and `this.method=something`, define class that `implements` the interface, and define the extension methods in this class.

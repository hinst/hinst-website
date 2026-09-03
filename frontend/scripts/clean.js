for (const directory of ['.parcel-cache', 'compiled', 'dist'])
	require('node:fs').rmSync(directory, { recursive: true, force: true });

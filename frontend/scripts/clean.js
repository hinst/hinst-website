for (const directory in ['.parcel-cache', 'compiled', 'dist'])
	require('node:fs').rmSync(directory, { recursive: true, force: true });

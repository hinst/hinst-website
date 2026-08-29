for (const directory in ['.parcel-cache', 'compiled', 'dist'])
	require('fs').rmSync(directory, { recursive: true, force: true });

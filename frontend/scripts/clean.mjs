import fs from 'node:fs';

for (const directory of ['.parcel-cache', 'compiled', 'dist'])
	fs.rmSync(directory, { recursive: true, force: true });

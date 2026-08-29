require('fs').rmSync('.parcel-cache', { recursive: true, force: true });
require('fs').rmSync('compiled', { recursive: true, force: true });
require('fs').rmSync('dist', { recursive: true, force: true });

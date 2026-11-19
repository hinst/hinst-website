// Build the hinst-website Docker image (linux/arm64)

import { execSync } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

const commands = [
	'docker run --rm --privileged multiarch/qemu-user-static --reset -p yes',
	'docker build -t 192.168.178.29:5000/hinst/hinst-website --progress=plain --platform=linux/arm64 .',
];

for (const command of commands) {
	console.log(command);
	execSync(command, { stdio: 'inherit', cwd: repoRoot });
}

// Push the hinst-website Docker image

import { execSync } from 'node:child_process';

const command = 'docker push 192.168.178.29:5000/hinst/hinst-website';

console.log(command);
execSync(command, { stdio: 'inherit' });

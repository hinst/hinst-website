// Build frontend UI to be displayed at GitHub.io pages

import { execSync } from 'child_process';
import fs from 'fs';

const subDirectory = '/dynamic';
const targetDirectory = '../../hinst.github.io' + subDirectory;
const API_URL = 'https://orangepizero2w-1.taile07783.ts.net/hinst-website/api';

fs.readdirSync(targetDirectory).forEach(file => {
	const filePath = targetDirectory + '/' + file;
	if (fs.lstatSync(filePath).isFile())
		fs.unlinkSync(filePath);
});

const command = 'npm run build -- --no-cache ' +
	'--dist-dir=' + targetDirectory +
	' --public-url=' + subDirectory;
console.log(command);
execSync(command, { stdio: 'inherit', env: { API_URL } });
const googleFileName = 'googled13030e7b9eaa45a.html';
fs.copyFileSync('./' + googleFileName, targetDirectory + '/' + googleFileName);

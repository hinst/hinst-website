import { execSync } from 'child_process';
import fs from 'fs';

const targetDirectory = '../../hinst.github.io';
const API_URL = 'https://orangepizero2w.taile07783.ts.net/hinst-website/api';

fs.readdirSync(targetDirectory).forEach(file => {
	const filePath = targetDirectory + '/' + file;
	if (fs.lstatSync(filePath).isFile())
		fs.unlinkSync(filePath);
});

execSync('npm run build -- --no-cache --dist-dir=' + targetDirectory, { stdio: 'inherit', env: { API_URL } });
const googleFileName = 'googled13030e7b9eaa45a.html';
fs.copyFileSync('./' + googleFileName, targetDirectory + '/' + googleFileName);

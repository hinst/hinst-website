import { execSync } from 'child_process';
import fs from 'fs';

const targetDirectory = '../../hinst.github.io';
const API_URL = 'https://orangepizero2w.tail46746a.ts.net/hinst-website/api';

fs.readdirSync(targetDirectory).forEach(file => {
	const filePath = targetDirectory + '/' + file;
	if (fs.lstatSync(filePath).isFile())
		fs.unlinkSync(filePath);
});

execSync('npm run build -- --no-cache --dist-dir=' + targetDirectory, { stdio: 'inherit', env: { API_URL } });

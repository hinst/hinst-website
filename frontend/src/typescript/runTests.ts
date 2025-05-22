import 'source-map-support/register';
import fs from 'fs';

function runTests(directory: string) {
	fs.readdirSync(directory).forEach((file) => {
		if (fs.lstatSync(directory + '/' + file).isDirectory()) runTests(directory + '/' + file);
		else if (file.endsWith('.test.js')) {
			console.log(`Running ${file}`);
			require('./' + file);
		}
	});
}

runTests('compiled');

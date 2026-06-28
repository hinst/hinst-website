import { Requests } from './requests.mjs';

async function main() {
	const requests = new Requests();
	const responseSizes = await requests.all();
}

main();

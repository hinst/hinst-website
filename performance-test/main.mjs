import { Requests } from "./requests.mjs";

async function main() {
	const requests = new Requests();
	await requests.all();
}

main();

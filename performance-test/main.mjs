//@ts-check
import { Requests } from './requests.mjs';

async function main() {
	const requests = new Requests();
	const initialSizes = await requests.all();
	console.log('Initial sizes:', initialSizes);

	const N = 10;
	const workers = [];

	for (let i = 0; i < N; i++) {
		workers.push(
			(async () => {
				while (true) {
					const sizes = await requests.all();
					for (let j = 0; j < sizes.length; j++) {
						if (sizes[j] !== initialSizes[j]) {
							throw new Error(
								`Size mismatch at index ${j}: expected ${initialSizes[j]}, got ${sizes[j]}`
							);
						}
					}
				}
			})()
		);
	}

	await Promise.all(workers);
}

main();

//@ts-check
import { Requests } from './requests.mjs';

const WORKER_COUNT = 30;

/**
	@param {number[]} initialSizes
	@param {number[]} sizes
*/
function assertSizes(initialSizes, sizes) {
	for (let i = 0; i < sizes.length; i++)
		if (sizes[i] !== initialSizes[i])
			throw new Error(
				`Size mismatch at index ${i}: expected ${initialSizes[i]}, got ${sizes[i]}`
			);
}

let requestCount = 0;
let totalRequestCount = 0;

/**
	@param {Requests} requests
	@param {number[]} initialSizes
*/
async function worker(requests, initialSizes) {
	while (true) {
		const sizes = await requests.all();
		assertSizes(sizes, initialSizes);
		requestCount++;
		totalRequestCount++;
	}
}

async function main() {
	const requests = new Requests();
	const initialSizes = await requests.all();
	console.log('Initial sizes:', initialSizes);

	const startTime = Date.now();
	setInterval(() => {
		const elapsed = (Date.now() - startTime) / 1000;
		const avgThroughput = (totalRequestCount / elapsed).toFixed(1);
		console.log(
			`Throughput: ${requestCount} requests.all() per second (avg: ${avgThroughput})`
		);
		requestCount = 0;
	}, 2000);

	const workers = [];
	for (let i = 0; i < WORKER_COUNT; i++) workers.push(worker(requests, initialSizes));
	await Promise.all(workers);
}

main();

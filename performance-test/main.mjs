//@ts-check
import { Requests } from './requests.mjs';

const WORKER_COUNT = 10;
const REPORT_INTERVAL_SECONDS = 10;

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

let blogPostCount = 0;
let totalBlogPostCount = 0;

/**
	@param {Requests} requests
	@param {number[]} initialSizes
*/
async function worker(requests, initialSizes) {
	while (true) {
		const sizes = await requests.all();
		assertSizes(initialSizes, sizes);
		const posts = requests.dateTimes.length;
		blogPostCount += posts;
		totalBlogPostCount += posts;
	}
}

async function main() {
	const requests = new Requests();
	console.time('Initializing...');
	const initialSizes = await requests.all();
	console.timeEnd('Initializing...');
	const megabytes = (initialSizes.reduce((a, b) => a + b, 0) / (1024 * 1024)).toFixed(1);
	console.log(`Initial sizes: ${megabytes} MB, total blog posts: ${requests.dateTimes.length}`);

	const startTime = Date.now();
	setInterval(() => {
		const elapsed = (Date.now() - startTime) / 1000;
		const avgThroughput = (totalBlogPostCount / elapsed).toFixed(1);
		const throughputPerSecond = (blogPostCount / REPORT_INTERVAL_SECONDS).toFixed(1);
		console.log(
			`Throughput: ${throughputPerSecond} blog posts per second (avg: ${avgThroughput})`
		);
		blogPostCount = 0;
	}, REPORT_INTERVAL_SECONDS * 1000);

	const workers = [];
	for (let i = 0; i < WORKER_COUNT; i++) workers.push(worker(new Requests(), initialSizes));
	await Promise.all(workers);
}

main();

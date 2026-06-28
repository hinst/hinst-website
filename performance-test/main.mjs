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
let requestCount = 0;
let totalRequestCount = 0;

/**
	@param {number[]} initialSizes
*/
async function worker(initialSizes) {
	while (true) {
		const requests = new Requests();
		const sizes = await requests.all();
		assertSizes(initialSizes, sizes);
		const posts = requests.dateTimes.length;
		blogPostCount += posts;
		totalBlogPostCount += posts;
		requestCount += requests.requestCount;
		totalRequestCount += requests.requestCount;
	}
}

async function main() {
	const requests = new Requests();
	console.time('Initializing...');
	const initialSizes = await requests.all();
	console.timeEnd('Initializing...');
	const megabytes = (initialSizes.reduce((a, b) => a + b, 0) / (1024 * 1024)).toFixed(1);
	console.log(
		`Initial sizes: ${megabytes} MB, total blog posts: ${requests.dateTimes.length}, total requests: ${requests.requestCount}`
	);

	const startTime = Date.now();
	setInterval(() => {
		const elapsed = (Date.now() - startTime) / 1000;
		const avgThroughput = (totalBlogPostCount / elapsed).toFixed(1);
		const throughputPerSecond = (blogPostCount / REPORT_INTERVAL_SECONDS).toFixed(1);
		const avgRequestsPerSecond = (totalRequestCount / elapsed).toFixed(1);
		const requestsPerSecond = (requestCount / REPORT_INTERVAL_SECONDS).toFixed(1);
		console.log(
			`Throughput: ${throughputPerSecond} blog posts per second (avg: ${avgThroughput}), ${requestsPerSecond} requests per second (avg: ${avgRequestsPerSecond})`
		);
		blogPostCount = 0;
		requestCount = 0;
	}, REPORT_INTERVAL_SECONDS * 1000);

	const workers = [];
	for (let i = 0; i < WORKER_COUNT; i++) workers.push(worker(initialSizes));
	await Promise.all(workers);
}

main();

import { RiddleEntry } from './riddle';

const webCounterUrl = 'http://localhost:8081/web-counter';

async function getRiddle(): Promise<RiddleEntry> {
	const url = webCounterUrl + '/riddle/generate';
	const response = await fetch(url);
	if (!response.ok)
		throw new Error(response.statusText);
	return Object.assign(new RiddleEntry(), await response.json());
}

async function getPrimeNumbers(): Promise<number[]> {
	const url = webCounterUrl + '/riddle/prime-numbers';
	const response = await fetch(url);
	if (!response.ok)
		throw new Error(response.statusText);
	return await response.json();
}

async function solveRiddle() {
	const riddle = await getRiddle();
	const primeNumbers = await getPrimeNumbers();
	console.time('solve');
	const indexes = riddle.solveFull(primeNumbers);
	console.timeEnd('solve');
	if (!indexes.length)
		throw new Error('Cannot solve riddle ' + riddle.id);
	return indexes;
}

export async function main() {
	const url = webCounterUrl + '/ping?url=' + encodeURIComponent(window.location.href);
	console.log(await solveRiddle());
}
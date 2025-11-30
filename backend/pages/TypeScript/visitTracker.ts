import { RiddleEntry } from './riddle';

const webCounterUrl = 'https://orangepizero2w-1.taile07783.ts.net/web-counter';

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
	const indexes = await riddle.solve(primeNumbers);
	if (!indexes.length)
		throw new Error('Cannot solve riddle ' + riddle.id);
	return riddle
}

export async function main() {
	const riddle = await solveRiddle();
	const url = webCounterUrl + '/ping' +
		'?url=' + encodeURIComponent(window.location.href) +
		'&riddleId=' + encodeURIComponent(riddle.id);
	await fetch(url, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(riddle.indexes),
	});
}

import { test } from 'node:test';
import assert from 'node:assert/strict';
import { RiddleSolver } from './riddle';

const primeNumbers = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 47];

function testSolve(chosenNumbers: number[]) {
	const steps = chosenNumbers.length;
	const limit = 1_000_000;
	const goal = chosenNumbers.reduce((acc, val) => acc * val, 1) % limit;
	const solver = new RiddleSolver(primeNumbers, goal, steps, limit);
	assert(solver.solve());
	assert.equal(solver.sequence.length, steps);
	const product = solver.sequence.reduce((acc, val) => (acc * val) % limit, 1);
	assert.equal(product, goal);

	const count = solver.count();
	console.log(chosenNumbers, count, solver.callCount);
}

test('RiddleSolver.solve simple', function () {
	const chosenNumbers = [2, 13, 23, 17];
	testSolve(chosenNumbers);
});

test('RiddleSolver.solve same', function () {
	const chosenNumbers = [41, 47, 41, 47];
	testSolve(chosenNumbers);
});

test('RiddleSolver.solve generated', function () {
	const testCount = 10;
	const steps = 4;
	for (let i = 0; i < testCount; i++) {
		const chosenNumbers = [];
		for (let j = 0; j < steps; j++) {
			const randomIndex = Math.floor(Math.random() * primeNumbers.length);
			chosenNumbers.push(primeNumbers[randomIndex]);
		}
		testSolve(chosenNumbers);
	}
});

export function main() {}

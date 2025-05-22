import { test } from 'node:test';
import assert from 'node:assert/strict';
import { RiddleSolver } from './riddle';

const primeNumbers = [2, 3, 5, 7, 11, 13, 17, 19, 23, 29];

test('RiddleSolver.solve simple', function () {
	const chosenNumbers = [2, 13, 23, 11];
	const steps = 4;
	const limit = 1000;
	const goal = chosenNumbers.reduce((acc, val) => acc * val, 1) % limit;
	const solver = new RiddleSolver(primeNumbers, goal, steps, limit);
	assert(solver.solve());
	assert(solver.sequence.length === steps);
	chosenNumbers.forEach((chosenNumber) => assert(solver.sequence.includes(chosenNumber)));
});

export function main() {}

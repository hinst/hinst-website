export class RiddleItem {
	id: number = 0;
	product: number = 0;
	steps: number = 0;
	limit: number = 0;

	solve(primeNumbers: number[]): number[] {
		const solver = new RiddleSolver(primeNumbers, this.product, this.steps, this.limit);
		if (solver.solve()) return solver.sequence;
		else return [];
	}
}

export class RiddleSolver {
	readonly sequence: number[];

	constructor(
		public readonly primeNumbers: number[],
		public readonly goal: number,
		public readonly steps: number,
		public readonly limit: number
	) {
		this.sequence = new Array(steps);
	}

	solve(step: number = 0, product: number = 1): boolean {
		if (step === this.steps) return product === this.goal;
		if (step > this.steps) throw new Error('Step out of bounds: ' + step);
		for (const primeNumber of this.primeNumbers) {
			this.sequence[step] = primeNumber;
			const newProduct = (product * primeNumber) % this.limit;
			if (this.solve(step + 1, newProduct)) return true;
		}
		return false;
	}
}

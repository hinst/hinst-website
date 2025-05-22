export class RiddleItem {
	id: number = 0;
	product: number = 0;
	steps: number = 0;

	solve(primeNumbers: number[]): number[] {
		const solver = new RiddleSolver(primeNumbers, this.product, this.steps);
		if (solver.solve()) return solver.sequence;
		else return [];
	}
}

export class RiddleSolver {
	readonly sequence: number[];

	constructor(
		public readonly primeNumbers: number[],
		public readonly goal: number,
		public readonly steps: number
	) {
		this.sequence = new Array(steps);
		for (let i = 0; i < steps; ++i) {
			this.sequence[i] = 1;
		}
	}

	solve(step: number = 0, product: number = 1): boolean {
		if (step === this.steps) return product === this.goal;
		if (step < this.steps)
			for (let i = 0; i < this.primeNumbers.length; ++i) {
				const primeNumber = this.primeNumbers[i];
				this.sequence[step] = primeNumber;
				const newProduct = product * primeNumber;
				if (this.solve(step + 1, newProduct)) return true;
			}
		return false;
	}
}

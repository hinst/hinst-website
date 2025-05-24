export class RiddleItem {
	id: number = 0;
	product: number = 0;
	steps: number = 0;
	limit: number = 1;
}

export class RiddleSolver {
	private readonly indexes: number[];
	callCount = 0;

	constructor(
		public readonly primeNumbers: number[],
		public readonly goal: number,
		public readonly steps: number,
		public readonly limit: number
	) {
		this.indexes = new Array(steps);
		this.indexes.fill(0);
	}

	get sequence(): number[] {
		return this.indexes.map((index) => this.primeNumbers[index]);
	}

	solve(step: number = 0, product: number = 1): boolean {
		this.callCount++;
		if (step === this.steps) return product === this.goal;
		if (step > this.steps) throw new Error('Step out of bounds: ' + step);
		const start = step === 0 ? 0 : this.indexes[step - 1];
		for (let i = start; i < this.primeNumbers.length; i++) {
			const primeNumber = this.primeNumbers[i];
			const newProduct = (product * primeNumber) % this.limit;
			this.indexes[step] = i;
			if (this.solve(step + 1, newProduct)) return true;
		}
		return false;
	}

	count(step: number = 0, product: number = 1): number {
		this.callCount++;
		if (step === this.steps) return product === this.goal ? 1 : 0;
		if (step > this.steps) throw new Error('Step out of bounds: ' + step);
		const start = step === 0 ? 0 : this.indexes[step - 1];
		let count = 0;
		for (let i = start; i < this.primeNumbers.length; i++) {
			const primeNumber = this.primeNumbers[i];
			const newProduct = (product * primeNumber) % this.limit;
			this.indexes[step] = i;
			count += this.count(step + 1, newProduct);
		}
		return count;
	}
}

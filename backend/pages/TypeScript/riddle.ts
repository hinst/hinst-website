export class RiddleEntry {
	constructor(
		public readonly id: number = 0,
		public readonly createdAt: string = new Date().toISOString(),
		public readonly product: number = 0,
		public readonly productLimit: number = 0,
		public readonly stepCount: number = 0,
	) {}

	solve(primeNumbers: number[]): number[] {
		const indexes: number[] = new Array(this.stepCount).fill(0);
		do {
			let product = 1;
			for (const index of indexes) {
				product *= primeNumbers[index] ?? 0;
				product %= this.productLimit;
			}
			if (product === this.product)
				return indexes;
		} while (next(indexes, primeNumbers.length))
		return [];
	}

	solveFull(primeNumbers: number[]): number[][] {
		const results: number[][] = [];
		const indexes: number[] = new Array(this.stepCount).fill(0);
		do {
			let product = 1;
			for (const index of indexes) {
				product *= primeNumbers[index] ?? 0;
				product %= this.productLimit;
			}
			if (product === this.product)
				results.push([...indexes]);
		} while (next(indexes, primeNumbers.length))
		return results;
	}
}

function next(sequence: number[], limit: number) {
	for (let i = 0; i < sequence.length; ++i) {
		const item = sequence[i] ?? 0;
		if (item < limit) {
			sequence[i] = item + 1;
			return true;
		}
		sequence[i] = 0;
	}
	return false;
}

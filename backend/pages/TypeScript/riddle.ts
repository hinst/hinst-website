import { sleep, Sleeper } from './sleep';

export class RiddleEntry {
	static readonly SLEEP_INTERVAL = 500_000;
	public indexes: number[] = [];

	constructor(
		public readonly id: number = 0,
		public readonly createdAt: string = new Date().toISOString(),
		public readonly product: number = 0,
		public readonly productLimit: number = 0,
		public readonly stepCount: number = 0,
	) {}

	async solve(primeNumbers: number[]): Promise<number[]> {
		this.indexes = new Array(this.stepCount).fill(0);
		const sleeper = new Sleeper(RiddleEntry.SLEEP_INTERVAL);
		do {
			let product = 1;
			for (const index of this.indexes) {
				product *= primeNumbers[index] ?? 0;
				product %= this.productLimit;
			}
			if (product === this.product)
				return this.indexes;
			if (sleeper.next())
				await sleep(0);
		} while (next(this.indexes, primeNumbers.length))
		return this.indexes = [];
		return [];
	}

	async solveFull(primeNumbers: number[]): Promise<number[][]> {
		const results: number[][] = [];
		const indexes: number[] = new Array(this.stepCount).fill(0);
		const sleeper = new Sleeper(RiddleEntry.SLEEP_INTERVAL);
		do {
			let product = 1;
			for (const index of indexes) {
				product *= primeNumbers[index] ?? 0;
				product %= this.productLimit;
			}
			if (product === this.product) {
				if (!this.indexes.length)
					this.indexes = [...indexes];
				results.push([...indexes]);
			}
			if (sleeper.next())
				await sleep(0);
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

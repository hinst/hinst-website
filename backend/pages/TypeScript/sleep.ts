export async function sleep(milliseconds: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, milliseconds));
}

export class Sleeper {
	private count = 0;

	constructor(public readonly intervalCount: number) {
	}

	next(): boolean {
		++this.count;
		if (this.count >= this.intervalCount) {
			this.count = 0;
			return true;
		}
		return false;
	}
}

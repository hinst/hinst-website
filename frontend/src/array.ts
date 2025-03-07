import lodash from 'lodash';

export function getPaddedChunks<T>(items: T[], size: number) {
	return lodash.chunk(items, size).map(posts => {
		const paddedItems: (T | undefined)[] = [...posts];
		while (paddedItems.length < size)
			paddedItems.push(undefined);
		return paddedItems;
	});
}

export function getPaddedArray<T>(items: T[], size: number) {
	const paddedItems: (T | undefined)[] = [...items];
	while (paddedItems.length % size !== 0)
		paddedItems.push(undefined);
	return paddedItems;
}


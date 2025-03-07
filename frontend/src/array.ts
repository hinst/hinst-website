import lodash from 'lodash';

export function getPaddedChunks<T>(posts: T[], size: number) {
	return lodash.chunk(posts, size).map(posts => {
		const paddedPosts: (T | undefined)[] = [...posts];
		while (paddedPosts.length < size)
			paddedPosts.push(undefined);
		return paddedPosts;
	});
}


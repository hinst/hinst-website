import { useEffect, useRef, useState } from 'react';
import { API_URL } from '../api';
import { NavLink, useParams } from 'react-router';
import { PostHeader } from './goalHeader';
import { compareStrings } from '../string';
import { getMonthName, parseMonthlyDate } from '../date';
import lodash from 'lodash';

class GroupedPosts {
	constructor(
		public readonly monthDate: string,
		public readonly posts: PostHeader[],
	) {}
}

export default function GoalPostListPanel() {
	const [isLoading, setIsLoading] = useState(0);
	const isLoadingRef = useRef(0);
	isLoadingRef.current = isLoading;

	const params = useParams();
	const id: string = params.id!;

	const [posts, setPosts] = useState(new Array<PostHeader>());

	async function loadPosts() {
		setIsLoading(isLoadingRef.current + 1);
		try {
			const response = await fetch(API_URL + '/goalPosts?id=' + encodeURIComponent(id));
			if (!response.ok)
				throw new Error(response.statusText);
			setPosts(await response.json());
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(() => { loadPosts() }, []);

	function getMonthlyPosts() {
		const sortedPosts = [...posts].sort((a, b) => compareStrings(a.date, b.date));
		const groups = new Map<string, PostHeader[]>();
		sortedPosts.forEach(post => {
			const monthDate = post.date.substring(0, '2025-03'.length);
			let group = groups.get(monthDate);
			if (!group) {
				group = new Array<PostHeader>();
				groups.set(monthDate, group);
			}
			group.push(post);
		});
		return Array.from(groups.entries())
			.map(([monthDate, posts]) => new GroupedPosts(monthDate, posts))
			.sort((a, b) => -compareStrings(a.monthDate, b.monthDate));
	}

	function getPaddedChunks<T>(posts: T[], size: number) {
		return lodash.chunk(posts, size).map(posts => {
			const paddedPosts: (T | undefined)[] = [...posts];
			while (paddedPosts.length < size)
				paddedPosts.push(undefined);
			return paddedPosts;
		});
	}

	return <div style={{display: 'flex', flexWrap: 'wrap', gap: 20}}>
		{ isLoading ? <div className='ms-loading'></div> : undefined }
		{getMonthlyPosts().map(group =>
			<div key={group.monthDate}
				className='ms-card ms-border'
				style={{display: 'flex', flexDirection: 'column', margin: 0, width: 'fit-content'}}
			>
				<div className='ms-card-title'>
					{group.monthDate.slice(0, '2025'.length)} &bull; {getMonthName(parseMonthlyDate(group.monthDate))}
				</div>
				<div style={{display: 'flex', flexDirection: 'column', gap: 10}}>
					{getPaddedChunks(group.posts, 3).map(posts =>
						<div style={{display: 'flex', gap: 10, flexWrap: 'wrap'}}>
							{posts.map(post =>
								post
									? <NavLink
										to={'/personal-goals/posts/' + post.id}
										key={post.id}
										className='ms-btn ms-primary ms-outline'
										style={{fontFamily: 'monospace'}}
									>
										{post.date.slice(8, '2025-03-06'.length)}
									</NavLink>
									: <NavLink
										to={'/personal-goals/posts/'}
										className='ms-btn ms-primary ms-outline'
										style={{fontFamily: 'monospace', visibility: 'hidden'}}
									>
										{'2025-03-06'.slice(8, '2025-03-06'.length)}
									</NavLink>
							)}
						</div>
					)}
				</div>
			</div>
		)}
	</div>;
}
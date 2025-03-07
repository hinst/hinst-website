import { useEffect, useRef, useState } from 'react';
import { API_URL } from '../api';
import { NavLink, useParams } from 'react-router';
import { PostHeader } from './goalHeader';
import { compareStrings } from '../string';
import { getMonthName, parseMonthlyDate } from '../date';
import { getPaddedChunks } from '../array';

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

	return <div style={{display: 'flex', flexWrap: 'wrap', gap: 20}}>
		{ isLoading ? <div className='ms-loading'></div> : undefined }
		{getPaddedChunks(getMonthlyPosts(), 6).map(rowGroup =>
			<div style={{display: 'flex', gap: 10, flexWrap: 'wrap'}}>
				{rowGroup.map(group =>
					group
					?
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
												? <NavLinkDay date={post.date} id={post.id}/>
												: <NavLinkDay date='2025-03-01' id={''}/>
										)}
									</div>
								)}
							</div>
						</div>
					: undefined
				)}
			</div>
		)}
	</div>;
}

function NavLinkDay(props: {date: string, id: string}) {
	return <NavLink
		to={'/personal-goals/posts/' + props.id}
		key={props.id}
		className='ms-btn ms-primary ms-outline'
		style={{fontFamily: 'monospace', visibility: props.id === '' ? 'hidden' : 'visible'}}
	>
		{props.date.slice('2025-03-'.length, '2025-03-06'.length)}
	</NavLink>;
}
import { useEffect, useRef, useState } from 'react';
import lodash from 'lodash';
import { API_URL } from '../api';
import { NavLink, useParams } from 'react-router';
import { PostHeader } from './goalHeader';
import { compareStrings } from '../string';
import { getMonthName, parseMonthlyDate } from '../date';
import { getPaddedArray, getPaddedChunks } from '../array';
import { Calendar } from 'react-feather';
import { createRandomId } from '../react';

const ROWS_PER_MONTH = 3;

class MonthlyPosts {
	constructor(
		public readonly monthDate: string,
		public readonly posts: PostHeader[],
	) {}
}

class YearlyPosts {
	constructor(
		public readonly year: string,
		public readonly monthlyPosts: MonthlyPosts[],
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

	function getSortedPosts() {
		return [...posts].sort((a, b) => compareStrings(a.date, b.date));
	}

	function getMonthlyPosts() {
		const posts = getSortedPosts();
		const groups = lodash.groupBy(posts, post => post.date.substring(0, '2025-03'.length));
		return Array.from(Object.entries(groups))
			.map(([monthDate, posts]) => new MonthlyPosts(monthDate, posts))
			.sort((a, b) => -compareStrings(a.monthDate, b.monthDate));
	}

	function getYearlyPosts() {
		const monthlyPosts = getMonthlyPosts();
		const groups = lodash.groupBy(monthlyPosts, monthlyPost => monthlyPost.monthDate.substring(0, '2025'.length));
		return Array.from(Object.entries(groups))
			.map(([year, monthlyPosts]) => new YearlyPosts(year, monthlyPosts))
			.sort((a, b) => -compareStrings(a.year, b.year));
	}

	return <div>
		{ isLoading ? <div className='ms-loading'></div> : undefined }
		{getYearlyPosts().map(yearGroup =>
			<div key={yearGroup.year}>
				<div

					style={{display: 'flex', gap: 10, alignItems: 'center'}}
				>
					<Calendar/>
					<div style={{whiteSpace: 'nowrap'}}>
						Year {yearGroup.year}
					</div>
					<hr/>
				</div>
				{yearGroup.monthlyPosts.map(group =>
					<div
						key={group.monthDate}
						className='ms-card ms-border'
						style={{
							width: 'fit-content',
							verticalAlign: 'top',
							marginTop: 0,
							marginRight: 15,
						}}
					>
						<div className='ms-card-title'>
							{getMonthName(parseMonthlyDate(group.monthDate))},&nbsp;
							{group.monthDate.slice(0, '2025'.length)}
						</div>
						<DaysOfMonth posts={group.posts}/>
					</div>
				)}
				<div/>
			</div>
		)}
	</div>;
}

function DaysOfMonth(props: {posts: PostHeader[]}) {
	return <div style={{display: 'flex', flexDirection: 'column', gap: 10}}>
		{getPaddedChunks(props.posts, ROWS_PER_MONTH).map((posts) =>
			<div
				key={posts.map(post => post?.id).join(' ')}
				style={{display: 'flex', gap: 10}}
			>
				{posts.map(post =>
					<div key={post?.id || createRandomId()}>
						{post
							? <NavLinkDay date={post.date} id={post.id}/>
							: <NavLinkDay date='2025-03-01' id={''}/>
						}
					</div>
				)}
			</div>
		)}
	</div>;
}

function NavLinkDay(props: {date: string, id: string}) {
	return <NavLink
		to={'/personal-goals/posts/' + props.id}
		className='ms-btn ms-primary ms-outline'
		style={{fontFamily: 'monospace', visibility: props.id === '' ? 'hidden' : 'visible'}}
	>
		{props.date.slice('2025-03-'.length, '2025-03-06'.length)}
	</NavLink>;
}
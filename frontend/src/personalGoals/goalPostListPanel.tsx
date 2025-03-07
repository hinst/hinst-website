import { useEffect, useRef, useState } from 'react';
import lodash from 'lodash';
import { API_URL } from '../api';
import { NavLink, useParams } from 'react-router';
import { PostHeader } from './goalHeader';
import { compareStrings } from '../string';
import { getMonthName, parseMonthlyDate } from '../date';
import { getPaddedArray } from '../array';

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
		<div>
			{getMonthlyPosts().map(group =>
				<div key={group.monthDate}
					className='ms-card ms-border'
					style={{
						display: 'inline-block',
						width: 'fit-content',
						verticalAlign: 'top',
						marginTop: 0,
						marginRight: 15,
						paddingBottom: 5,
						paddingRight: 5,
					}}
				>
					<div className='ms-card-title' style={{display: 'inline-block'}}>
						{group.monthDate.slice(0, '2025'.length)} &bull; {getMonthName(parseMonthlyDate(group.monthDate))}
					</div>
					<br/>
					{getPaddedArray(group.posts, ROWS_PER_MONTH).map((post, index) => [
							<div style={{display: 'inline-block', marginRight: 10, marginBottom: 10}}>
								{post
									? <NavLinkDay date={post.date} id={post.id}/>
									: <NavLinkDay date='2025-03-01' id={''}/>
								}
							</div>,
							(index + 1) % ROWS_PER_MONTH === 0 ? <br/> : undefined
						]
					)}
				</div>
			)}
		</div>
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
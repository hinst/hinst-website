import lodash from 'lodash';
import { NavLink } from 'react-router';
import { PostHeader } from './goalHeader';
import { compareStrings } from '../string';
import { getMonthName, parseMonthlyDate } from '../date';
import { getPaddedChunks } from '../array';
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

export default function GoalCalendar(props: {posts: PostHeader[], activePostDate: string}) {
	function getSortedPosts() {
		return [...props.posts].sort((a, b) => compareStrings(a.date, b.date));
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
		{getYearlyPosts().map(yearGroup =>
			<div key={yearGroup.year}>
				<div
					style={{display: 'flex', gap: 10, alignItems: 'center', marginBottom: 10}}
				>
					<Calendar/>
					<div style={{whiteSpace: 'nowrap'}}>
						Year {yearGroup.year}
					</div>
				</div>
				{yearGroup.monthlyPosts.map(group =>
					<div
						key={group.monthDate}
						className='ms-card ms-border'
						style={{
							width: 'fit-content',
							verticalAlign: 'top',
							marginTop: 0,
						}}
					>
						<div className='ms-card-title'>
							{getMonthName(parseMonthlyDate(group.monthDate))},&nbsp;
							{group.monthDate.slice(0, '2025'.length)}
						</div>
						<DaysOfMonth posts={group.posts} activePostDate={props.activePostDate} />
					</div>
				)}
				<div/>
			</div>
		)}
	</div>;
}

function DaysOfMonth(props: {posts: PostHeader[], activePostDate: string}) {
	return <div style={{display: 'flex', flexDirection: 'column', gap: 10}}>
		{getPaddedChunks(props.posts, ROWS_PER_MONTH).map((posts) =>
			<div
				key={posts.map(post => post?.id).join(' ')}
				style={{display: 'flex', gap: 10}}
			>
				<DaysOfMonthRow posts={posts} activePostDate={props.activePostDate}/>
			</div>
		)}
	</div>;
}

function DaysOfMonthRow(props: {
	posts: (PostHeader | undefined)[],
	activePostDate: string,
}) {
	return props.posts.map(post =>
		<div key={post?.id || createRandomId()}>
			{post
				? <NavLinkDay
					date={post.date}
					id={post.id}
					goalId={post.obj_id}
					isActive={props.activePostDate === post.date}
					/>
				: <NavLinkDay
					date='2025-03-01'
					id={''}
					goalId={''}
					isActive={false}
					/>
			}
		</div>
	);
}

function NavLinkDay(props: {date: string, id: string, goalId: string, isActive: boolean}) {
	let url = '';
	if (props.id !== '')
		url = '/personal-goals/' + encodeURIComponent(props.goalId) +
			'?activePostDate=' + encodeURIComponent(props.date);
	return <NavLink
		to={url}
		className={'ms-btn ms-outline ms-primary ' + (props.isActive ? 'ms-btn-active' : '')}
		style={{fontFamily: 'monospace', visibility: props.id === '' ? 'hidden' : 'visible'}}
	>
		{props.date.slice('2025-03-'.length, '2025-03-06'.length)}
	</NavLink>;
}
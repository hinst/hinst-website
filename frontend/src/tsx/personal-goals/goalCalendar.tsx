import lodash from 'lodash';
import { NavLink } from 'react-router';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { compareStrings } from 'src/typescript/string';
import { getMonthName, parseMonthlyDate } from 'src/typescript/date';
import { getPaddedChunks } from 'src/typescript/array';
import { Calendar } from 'react-feather';
import { createRandomId } from 'src/typescript/react';
import { DateTime } from 'luxon';

const ROWS_PER_MONTH = 3;

class MonthlyPosts {
	constructor(
		public readonly monthDate: string,
		public readonly posts: GoalPostRecord[]
	) {}
}

class YearlyPosts {
	constructor(
		public readonly year: string,
		public readonly monthlyPosts: MonthlyPosts[]
	) {}
}

export default function GoalCalendar(props: { posts: GoalPostRecord[]; activePostDate: number }) {
	function getSortedPosts() {
		return [...props.posts].sort((a, b) => a.dateTime - b.dateTime);
	}

	function getMonthlyPosts() {
		const posts = getSortedPosts();
		const groups = lodash.groupBy(posts, (post) => post.yearAndMonthText);
		return Array.from(Object.entries(groups))
			.map(([monthDate, posts]) => new MonthlyPosts(monthDate, posts))
			.sort((a, b) => -compareStrings(a.monthDate, b.monthDate));
	}

	function getYearlyPosts() {
		const monthlyPosts = getMonthlyPosts();
		const groups = lodash.groupBy(monthlyPosts, (monthlyPost) =>
			monthlyPost.monthDate.substring(0, '2025'.length)
		);
		return Array.from(Object.entries(groups))
			.map(([year, monthlyPosts]) => new YearlyPosts(year, monthlyPosts))
			.sort((a, b) => -compareStrings(a.year, b.year));
	}

	return (
		<div>
			{getYearlyPosts().map((yearGroup) => (
				<div key={yearGroup.year}>
					<div
						style={{ display: 'flex', gap: 10, alignItems: 'center', marginBottom: 10 }}
					>
						<Calendar />
						<div style={{ whiteSpace: 'nowrap' }}>Year {yearGroup.year}</div>
					</div>
					{yearGroup.monthlyPosts.map((group) => (
						<div
							key={group.monthDate}
							className='ms-card ms-border'
							style={{
								width: 'fit-content',
								verticalAlign: 'top',
								marginTop: 0
							}}
						>
							<div className='ms-card-title'>
								{getMonthName(parseMonthlyDate(group.monthDate))},&nbsp;
								{group.monthDate.slice(0, '2025'.length)}
							</div>
							<DaysOfMonth
								posts={group.posts}
								activePostDate={props.activePostDate}
							/>
						</div>
					))}
					<div />
				</div>
			))}
		</div>
	);
}

function DaysOfMonth(props: { posts: GoalPostRecord[]; activePostDate: number }) {
	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			{getPaddedChunks(props.posts, ROWS_PER_MONTH).map((posts) => (
				<div
					key={posts.map((post) => post?.dateTime).join(' ')}
					style={{ display: 'flex', gap: 10 }}
				>
					<DaysOfMonthRow posts={posts} activePostDate={props.activePostDate} />
				</div>
			))}
		</div>
	);
}

function DaysOfMonthRow(props: { posts: (GoalPostRecord | undefined)[]; activePostDate: number }) {
	return props.posts.map((post) => (
		<div key={post?.dateTime || createRandomId()}>
			{post ? (
				<NavLinkDay
					date={post.dateTime}
					goalId={post.goalId}
					isActive={props.activePostDate === post.dateTime}
					isPublic={post.isPublic}
				/>
			) : (
				<NavLinkDay date={0} goalId={0} isActive={false} isPublic={false} />
			)}
		</div>
	));
}

function NavLinkDay(props: { date: number; goalId: number; isActive: boolean; isPublic: boolean }) {
	let url = '';
	if (props.date !== 0)
		url =
			'/personal-goals/' +
			encodeURIComponent(props.goalId) +
			'?activePostDate=' +
			encodeURIComponent(props.date);
	const classNames = ['ms-btn', 'ms-outline'];
	classNames.push(props.isPublic ? 'ms-primary' : 'ms-secondary ms-text-secondary');
	if (props.isActive) classNames.push('ms-btn-active');
	const dayText = DateTime.fromMillis(props.date * 1000).toFormat('dd');
	return (
		<NavLink
			to={url}
			className={classNames.join(' ')}
			style={{ fontFamily: 'monospace', visibility: props.date === 0 ? 'hidden' : 'visible' }}
		>
			{dayText}
		</NavLink>
	);
}

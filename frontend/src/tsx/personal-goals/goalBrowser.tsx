import { DateTime } from 'luxon';
import { useContext, useEffect, useRef, useState } from 'react';
import { useParams, useSearchParams } from 'react-router';
import { apiClient } from 'src/typescript/apiClient';
import { AppContext } from 'src/typescript/appContext';
import { PageTitle } from 'src/typescript/pageTitle';
import type { GoalPostHeaderEx } from 'src/typescript/rest_objects/goalPostHeaderEx';
import { requireString } from 'src/typescript/string';
import GoalBrowserNarrow from './goalBrowser.narrow';
import GoalBrowserWide from './goalBrowser.wide';
import GoalCalendar from './goalCalendar';
import GoalPostPanel from './goalPostPanel';

export default function GoalBrowser() {
	const context = useContext(AppContext);
	const params = useParams();
	const goalId: string = requireString(params.id);
	const [searchParams, setSearchParams] = useSearchParams();
	const activePostDate = searchParams.get('activePostDate') || '';

	const [goalTitle, setGoalTitle] = useState('');
	const [reloadGoalCalendar, setReloadGoalCalendar] = useState(0);

	const [isLoading, setIsLoading] = useState(0);
	const isLoadingRef = useRef(0);
	isLoadingRef.current = isLoading;

	const [posts, setPosts] = useState([] as GoalPostHeaderEx[]);

	const isFullMode = context.windowWidth >= 700;

	const [calendarVisible, setCalendarVisible] = useState(isFullMode);
	const [calendarTransition, setCalendarTransition] = useState('');

	function receivePosts(posts: GoalPostHeaderEx[]) {
		if (posts.length && !activePostDate) {
			const newActivePostDate = posts[0].dateTime;
			setSearchParams({ activePostDate: '' + newActivePostDate }, { replace: true });
		}
	}

	async function loadGoal(goalId: string) {
		const goalHeader = await apiClient.getGoal(parseInt(goalId, 10));
		setGoalTitle(goalHeader.getTitle(context.currentLanguage));
	}

	useEffect(() => {
		const _ = loadGoal(goalId);
	}, [goalId]);

	async function loadPosts() {
		setIsLoading(isLoadingRef.current + 1);
		try {
			const posts = await apiClient.getGoalPosts(parseInt(goalId, 10) || 0);
			setPosts(posts);
			receivePosts(posts);
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(() => {
		const _ = loadPosts();
	}, [goalId, reloadGoalCalendar]);

	useEffect(() => {
		if (activePostDate) setCalendarVisible(false);
		setTimeout(() => setCalendarTransition('transform 0.3s'));
	}, [activePostDate]);

	useEffect(() => {
		const activePostDateTime = DateTime.fromMillis(parseInt(activePostDate, 10) * 1000);
		const activePostDateTimeText = activePostDateTime.isValid
			? context.isAdminMode
				? activePostDateTime.toLocaleString({ dateStyle: 'short', timeStyle: 'short' })
				: activePostDateTime.toLocaleString({ dateStyle: 'short' })
			: '';
		context.setPageTitle(new PageTitle(goalTitle, activePostDateTimeText));
	}, [goalTitle, activePostDate, context.isAdminMode]);

	function getGoalCalendarPanel() {
		return (
			<div>
				{isLoading ? <div className='ms-loading' /> : undefined}
				<GoalCalendar posts={posts} activePostDate={parseInt(activePostDate, 10) || 0} />
			</div>
		);
	}

	function getGoalPostPanel() {
		return (
			<GoalPostPanel
				goalId={parseInt(goalId, 10)}
				postDate={parseInt(activePostDate, 10)}
				onChange={() => setReloadGoalCalendar(Math.random())}
			/>
		);
	}

	return isFullMode ? (
		<GoalBrowserWide
			activePostDate={activePostDate}
			getGoalCalendarPanel={getGoalCalendarPanel}
			getGoalPostPanel={getGoalPostPanel}
		/>
	) : (
		<GoalBrowserNarrow
			activePostDate={activePostDate}
			calendarVisible={calendarVisible}
			setCalendarVisible={setCalendarVisible}
			calendarTransition={calendarTransition}
			getGoalCalendarPanel={getGoalCalendarPanel}
			getGoalPostPanel={getGoalPostPanel}
		/>
	);
}

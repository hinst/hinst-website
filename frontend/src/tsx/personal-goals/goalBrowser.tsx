import { useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { GoalPostHeaderEx } from 'src/typescript/rest_objects/goalPostHeaderEx';
import { useContext, useEffect, useState } from 'react';
import { AppContext } from 'src/tsx/context';
import Cookie from 'js-cookie';
import GoalPostPanel from './goalPostPanel';
import GoalBrowserNarrow from './goalBrowser.narrow';
import { DateTime } from 'luxon';
import { apiClient } from 'src/typescript/apiClient';

const ARTICLE_PADDING = 20;
const ARTICLE_WIDTH = 1000 + ARTICLE_PADDING * 2;
const STRIPES_BACKGROUND = `repeating-linear-gradient(
	45deg,
	rgba(var(--main-bg), 1),
	rgba(var(--main-bg), 1) 10px,
	rgba(var(--light-bg-color), 1) 10px,
	rgba(var(--light-bg-color), 1) 20px
)`;
const STRIPES_MIN_WIDTH = 100;

export default function GoalBrowser(props: { setPageTitle: (title: string) => void }) {
	const context = useContext(AppContext);
	const params = useParams();
	const goalId: string = params.id!;
	const [searchParams, setSearchParams] = useSearchParams();
	const activePostDate = searchParams.get('activePostDate') || '';
	const [articleContainerId] = useState('id-' + Math.random().toString(16));
	const [articleContainerWidth, setArticleContainerWidth] = useState(0);

	const [goalTitle, setGoalTitle] = useState('');
	const [reloadGoalCalendar, setReloadGoalCalendar] = useState(0);

	function isFullMode() {
		return context.windowWidth >= 700;
	}

	function isGoalManagerMode() {
		return Cookie.get('goalManagerMode') === '1';
	}

	const [calendarVisible, setCalendarVisible] = useState(isFullMode());
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
		loadGoal(goalId);
	}, [goalId]);

	useEffect(() => {
		if (activePostDate) setCalendarVisible(false);
		setTimeout(() => setCalendarTransition('transform 0.3s'));
	}, [activePostDate]);

	useEffect(() => {
		const dateTime = DateTime.fromMillis(parseInt(activePostDate) * 1000);
		const dateText = isGoalManagerMode()
			? dateTime.toLocaleString({ dateStyle: 'short', timeStyle: 'short' })
			: dateTime.toLocaleString({ dateStyle: 'short' });
		const components = [goalTitle, dateText].filter((s) => s.length);
		props.setPageTitle(components.join(' • '));
	}, [goalTitle, activePostDate]);

	useEffect(() => {
		const timer = setInterval(() => {
			const container = document.getElementById(articleContainerId);
			setArticleContainerWidth(container?.clientWidth || 0);
		});
		return () => clearInterval(timer);
	}, []);

	function getGoalCalendarPanel() {
		return (
			<GoalCalendarPanel
				id={goalId}
				receivePosts={receivePosts}
				activePostDate={parseInt(activePostDate) || 0}
				reload={reloadGoalCalendar}
			/>
		);
	}

	function getGoalPostPanel() {
		return (
			<GoalPostPanel
				goalId={parseInt(goalId)}
				postDate={parseInt(activePostDate)}
				goalManagerMode={isGoalManagerMode()}
				onChange={() => setReloadGoalCalendar(Math.random())}
			/>
		);
	}

	function getWideLayout() {
		return (
			<div
				style={{
					display: 'flex',
					gap: 20,
					minHeight: 0,
					height: '100%'
				}}
			>
				<div
					style={{
						display: 'flex',
						overflowY: 'auto',
						flexShrink: 0,
						flexBasis: 'fit-content'
					}}
				>
					{getGoalCalendarPanel()}
				</div>
				<div
					id={articleContainerId}
					style={{
						flexGrow: 1,
						justifyContent: 'center',
						display: 'flex',
						minHeight: 0,
						maxHeight: '100%',
						background:
							articleContainerWidth > ARTICLE_WIDTH + STRIPES_MIN_WIDTH
								? STRIPES_BACKGROUND
								: undefined
					}}
				>
					<div
						className='ms-bg-main'
						style={{
							paddingLeft: ARTICLE_PADDING,
							paddingRight: ARTICLE_PADDING,
							flexGrow: 1,
							maxWidth: ARTICLE_WIDTH,
							backgroundAttachment: 'fixed',
							minHeight: 0,
							overflowY: 'auto'
						}}
					>
						{activePostDate ? getGoalPostPanel() : undefined}
					</div>
				</div>
			</div>
		);
	}

	return isFullMode() ? (
		getWideLayout()
	) : (
		<GoalBrowserNarrow
			goalId={goalId}
			activePostDate={activePostDate}
			calendarVisible={calendarVisible}
			setCalendarVisible={setCalendarVisible}
			calendarTransition={calendarTransition}
			goalManagerMode={isGoalManagerMode()}
			receivePosts={receivePosts}
			reloadGoalCalendar={reloadGoalCalendar}
			onChangeGoalCalendar={() => setReloadGoalCalendar(Math.random())}
		/>
	);
}

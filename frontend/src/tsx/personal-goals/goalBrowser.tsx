import { useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { GoalRecord } from 'src/typescript/personal-goals/goalRecord';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { useContext, useEffect, useState } from 'react';
import { API_URL } from 'src/typescript/global';
import { translateGoalTitle } from 'src/typescript/personal-goals/goalInfo';
import { AppContext } from 'src/tsx/context';
import { Calendar } from 'react-feather';
import Cookie from 'js-cookie';
import GoalPostPanel from './goalPostPanel';
import { DateTime } from 'luxon';

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

	const [calendarEnabled, setCalendarEnabled] = useState(true);
	const [calendarTransition, setCalendarTransition] = useState('');

	function receivePosts(posts: GoalPostRecord[]) {
		if (posts.length && !activePostDate) {
			const newActivePostDate = posts[posts.length - 1].dateTime;
			setSearchParams({ activePostDate: '' + newActivePostDate }, { replace: true });
		}
	}

	async function loadGoal(goalId: string) {
		const response = await fetch(API_URL + '/goal?id=' + encodeURIComponent(goalId));
		if (response.ok) {
			const goalHeader: GoalRecord = await response.json();
			setGoalTitle(translateGoalTitle(context.currentLanguage, goalHeader.title));
		}
	}

	useEffect(() => {
		loadGoal(goalId);
	}, [goalId]);

	useEffect(() => {
		if (activePostDate) setCalendarEnabled(false);
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

	function getFloatingCalendarButton() {
		return (
			<div
				className='ms-bg-light ms-shape-circle'
				style={{
					position: 'absolute',
					width: 40,
					height: 40,
					bottom: 0,
					right: 0,
					zIndex: 2
				}}
			>
				<button
					className={
						'ms-btn ms-primary ms-rounded ms-box-shadow' +
						(calendarEnabled ? ' ms-btn-active' : '')
					}
					onClick={() => setCalendarEnabled(!calendarEnabled)}
					style={{
						margin: 0,
						width: 40,
						height: 40
					}}
				>
					<Calendar
						style={{
							position: 'absolute',
							left: '50%',
							top: '50%',
							transform: 'translate(-50%, -50%)'
						}}
					/>
				</button>
			</div>
		);
	}

	function getNarrowLayout() {
		return (
			<div
				style={{
					position: 'relative',
					display: 'flex',
					minHeight: 0,
					height: '100%',
					maxHeight: '100%',
					width: '100%',
					maxWidth: '100%',
					overflowY: 'hidden'
				}}
			>
				{getFloatingCalendarButton()}
				<div
					className='ms-bg-light ms-shape-round ms-box-shadow ms-border-main'
					style={{
						display: calendarEnabled ? 'block' : 'block',
						position: 'absolute',
						zIndex: 1,
						overflowY: 'auto',
						maxHeight: '100%',
						padding: 8,
						borderWidth: 1,
						borderStyle: 'solid',
						transform: calendarEnabled ? 'translate(0,0)' : 'translate(-100%, 0)',
						transition: calendarTransition
					}}
				>
					{getGoalCalendarPanel()}
				</div>
				<div
					onClick={() => setCalendarEnabled(false)}
					style={{
						display: 'flex',
						overflowY: 'auto',
						flexGrow: 1
					}}
				>
					<div style={{ maxWidth: ARTICLE_WIDTH, width: '100%' }}>
						{activePostDate ? getGoalPostPanel() : undefined}
					</div>
				</div>
			</div>
		);
	}

	return isFullMode() ? getWideLayout() : getNarrowLayout();
}

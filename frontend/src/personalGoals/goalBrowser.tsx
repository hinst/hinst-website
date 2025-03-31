import { useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { GoalHeader, PostHeader } from './goalHeader';
import { useContext, useEffect, useState } from 'react';
import { API_URL } from '../global';
import { translateGoalTitle } from './goalInfo';
import { DisplayWidthContext, LanguageContext } from '../context';
import { Calendar } from 'react-feather';
import Cookie from 'js-cookie';
import GoalPostPanel from './goalPostPanel';

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

export default function GoalBrowser(props: {
	setPageTitle: (title: string) => void
}) {
	const currentLanguage = useContext(LanguageContext);
	const displayWidth = useContext(DisplayWidthContext);
	const params = useParams();
	const goalId: string = params.id!;
	const [searchParams, setSearchParams] = useSearchParams();
	const activePostDate = searchParams.get('activePostDate') || '';
	const [articleContainerId] = useState('id-' + Math.random().toString(16));
	const [articleContainerWidth, setArticleContainerWidth] = useState(0);

	const [goalTitle, setGoalTitle] = useState('');
	const [reloadGoalCalendar, setReloadGoalCalendar] = useState(0);

	function isFullMode() {
		return displayWidth >= 700;
	};

	function isGoalManagerMode() {
		return Cookie.get('goalManagerMode') === '1';
	}

	const [calendarEnabled, setCalendarEnabled] = useState(true);
	const [calendarTransition, setCalendarTransition] = useState('');

	function receivePosts(posts: PostHeader[]) {
		if (posts.length && !activePostDate) {
			const newActivePostDate = posts[posts.length - 1].date;
			setSearchParams({activePostDate: newActivePostDate});
		}
	};

	async function loadGoal() {
		const response = await fetch(API_URL + '/goal?id=' + encodeURIComponent(goalId));
		if (response.ok) {
			const goalHeader: GoalHeader = await response.json();
			setGoalTitle(translateGoalTitle(currentLanguage, goalHeader.title));
		}
	};

	useEffect(() => {
		loadGoal();
	}, [goalId]);

	useEffect(() => {
		if (activePostDate)
			setCalendarEnabled(false);
		setTimeout(() => setCalendarTransition('transform 0.3s'));
	}, [activePostDate]);

	useEffect(() => {
		const dateText = isGoalManagerMode()
			? activePostDate
			: activePostDate.slice(0, '2025-03-10'.length);
		const components = [goalTitle, dateText].filter(s => s.length);
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
		return <GoalCalendarPanel
			id={goalId}
			receivePosts={receivePosts}
			activePostDate={activePostDate}
			reload={reloadGoalCalendar}
		/>;
	};

	function getGoalPostPanel() {
		return <GoalPostPanel
			goalId={goalId}
			postDate={activePostDate}
			goalManagerMode={isGoalManagerMode()}
			onChange={() => setReloadGoalCalendar(Math.random())}
		/>;
	};

	function getWideLayout() {
		return <div style={{display: 'flex', gap: 20, minHeight: 0}}>
			<div style={{
				display: 'flex',
				overflowY: 'auto',
				paddingRight: 10,
				flexShrink: 0,
				flexBasis: 'fit-content',
			}}>
				{ getGoalCalendarPanel() }
			</div>
			<div
				id={articleContainerId}
				style={{
					flexGrow: 1,
					justifyContent: 'center',
					display: 'flex',
					minHeight: 0,
					maxHeight: '100%',
					background: articleContainerWidth > ARTICLE_WIDTH + STRIPES_MIN_WIDTH
						? STRIPES_BACKGROUND
						: undefined,
				}}
			>
				<div
					className='ms-bg-main'
					style={{
						paddingLeft: ARTICLE_PADDING,
						paddingRight: ARTICLE_PADDING,
						maxWidth: ARTICLE_WIDTH,
						backgroundAttachment: 'fixed',
						minHeight: 0,
						overflowY: 'auto',
					}}
				>
					{ activePostDate ? getGoalPostPanel() : undefined }
				</div>
			</div>
		</div>;
	};

	function getFloatingCalendarButton() {
		return <div
			className='ms-bg-light ms-shape-circle'
			style={{
				position: 'absolute',
				width: 40,
				height: 40,
				bottom: 0,
				right: 0,
				zIndex: 2,
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
					height: 40,
				}}
			>
				<Calendar
					style={{
						position: 'absolute',
						left: '50%',
						top: '50%',
						transform: 'translate(-50%, -50%)',
					}}
				/>
			</button>
		</div>;
	};

	function getNarrowLayout() {
		return <div
			style={{
				position: 'relative',
				display: 'flex',
				minHeight: 0,
				height: '100%',
				maxHeight: '100%',
				width: '100%',
				maxWidth: '100%',
				overflowY: 'hidden',
			}}
		>
			{ getFloatingCalendarButton() }
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
					transition: calendarTransition,
				}}
			>
				{ getGoalCalendarPanel() }
			</div>
			<div
				onClick={() => setCalendarEnabled(false)}
				style={{
					display: 'flex',
					overflowY: 'auto',
				}}
			>
				<div style={{maxWidth: ARTICLE_WIDTH}}>
					{ activePostDate ? getGoalPostPanel() : undefined }
				</div>
			</div>
		</div>;
	};

	return isFullMode() ? getWideLayout() : getNarrowLayout();
}
import { useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { GoalHeader, PostHeader } from './goalHeader';
import { useContext, useEffect, useState } from 'react';
import GoalPostView from './goalPostView';
import { API_URL } from '../global';
import { translateGoalTitle } from './goalInfo';
import { DisplayWidthContext, LanguageContext } from '../context';
import { Calendar } from 'react-feather';

export default function GoalBrowser(props: {
	setPageTitle: (title: string) => void
}) {
	const currentLanguage = useContext(LanguageContext);
	const displayWidth = useContext(DisplayWidthContext);
	const params = useParams();
	const goalId: string = params.id!;
	const [searchParams, setSearchParams] = useSearchParams();
	const activePostDate = searchParams.get('activePostDate') || '';

	function isFullMode() {
		return displayWidth >= 700;
	}

	const [calendarEnabled, setCalendarEnabled] = useState(true);

	function receivePosts(posts: PostHeader[]) {
		if (posts.length && !activePostDate) {
			const newActivePostDate = posts[posts.length - 1].date;
			setSearchParams({activePostDate: newActivePostDate});
		}
	}

	async function loadGoal() {
		const response = await fetch(API_URL + '/goal?id=' + encodeURIComponent(goalId));
		if (response.ok) {
			const goalHeader: GoalHeader = await response.json();
			const goalTitle = translateGoalTitle(currentLanguage, goalHeader.title);
			props.setPageTitle(goalTitle);
		}
	}

	useEffect(() => {
		loadGoal();
	}, [goalId]);

	useEffect(() => {
		if (activePostDate)
			setCalendarEnabled(false);
	}, [activePostDate]);

	function getWideLayout() {
		return <div style={{display: 'flex', gap: 20, minHeight: 0}}>
			<div style={{
				display: 'flex',
				overflowY: 'auto',
				paddingRight: 10,
				flexShrink: 0,
				flexBasis: 'fit-content',
			}}>
				<GoalCalendarPanel
					id={goalId}
					receivePosts={receivePosts}
					activePostDate={activePostDate}
				/>
			</div>
			<div style={{
				display: 'flex',
				overflowY: 'auto',
				flexGrow: 1
			}}>
				{ activePostDate
					? <GoalPostView
						goalId={goalId}
						postDate={activePostDate}
						style={{maxWidth: 1000}}
					/>
					: undefined }
			</div>
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
			<div
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
			</div>
			<div
				className='ms-bg-light ms-shape-round ms-box-shadow'
				style={{
					display: calendarEnabled ? 'block' : 'none',
					position: 'absolute',
					zIndex: 1,
					overflowY: 'auto',
					maxHeight: '100%',
					padding: 8,
				}}
			>
				<GoalCalendarPanel
					id={goalId}
					receivePosts={receivePosts}
					activePostDate={activePostDate}
				/>
			</div>
			<div style={{
				display: 'flex',
				overflowY: 'auto',
			}}>
				{ activePostDate
					? <GoalPostView
						goalId={goalId}
						postDate={activePostDate}
						style={{maxWidth: 1000}}
					/>
					: undefined }
			</div>
		</div>;
	};

	return isFullMode() ? getWideLayout() : getNarrowLayout();
}
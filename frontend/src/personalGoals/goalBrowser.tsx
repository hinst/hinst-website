import { useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { GoalHeader, PostHeader } from './goalHeader';
import { useContext, useEffect, useState } from 'react';
import GoalPostView from './goalPostView';
import { API_URL } from '../global';
import { translateGoalTitle } from './goalInfo';
import { DisplayWidthContext, LanguageContext } from '../context';

export default function GoalBrowser(props: {
	setPageTitle: (title: string) => void
}) {
	const currentLanguage = useContext(LanguageContext);
	const displayWidth = useContext(DisplayWidthContext);
	const params = useParams();
	const [searchParams, setSearchParams] = useSearchParams();
	const postId: string = params.id!;

	function isFullMode() {
		return displayWidth >= 700;
	}

	const [activePostDate, setActivePostDate] =
		useState<string>(searchParams.get('activePostDate') || '');
	const [calendarEnabled, setCalendarEnabled] = useState(true);

	function isCalendarVisible() {
		return isFullMode() || calendarEnabled;
	}

	function receivePosts(posts: PostHeader[]) {
		if (posts.length && !activePostDate) {
			const newActivePostDate = posts[posts.length - 1].date;
			setSearchParams({activePostDate: newActivePostDate});
			setActivePostDate(newActivePostDate);
		}
	}

	useEffect(() => {
		setActivePostDate(searchParams.get('activePostDate') || '');
	}, [searchParams]);

	async function loadGoal() {
		const response = await fetch(API_URL + '/goal?id=' + encodeURIComponent(postId));
		if (response.ok) {
			const goalHeader: GoalHeader = await response.json();
			const goalTitle = translateGoalTitle(currentLanguage, goalHeader.title);
			props.setPageTitle(goalTitle);
		}
	}

	useEffect(() => {
		loadGoal();
	}, [params]);

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
					id={postId}
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
						goalId={postId}
						postDate={activePostDate}
						style={{maxWidth: 1000}}
					/>
					: undefined }
			</div>
		</div>;
	};

	function getNarrowLayout() {
		return <div style={{
			position: 'relative',
			display: 'flex',
			minHeight: 0,
			height: '100%',
			maxHeight: '100%',
			width: '100%',
			maxWidth: '100%',
			overflowY: 'hidden',
		}}>
			<div
				className='ms-bg-light ms-shape-round ms-box-shadow'
				style={{
					position: 'absolute',
					zIndex: 1,
					overflowY: 'auto',
					maxHeight: '100%',
					padding: 8,
				}}
			>
				<GoalCalendarPanel
					id={postId}
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
						goalId={postId}
						postDate={activePostDate}
						style={{maxWidth: 1000}}
					/>
					: undefined }
			</div>
		</div>;
	};

	return isFullMode() ? getWideLayout() : getNarrowLayout();
}
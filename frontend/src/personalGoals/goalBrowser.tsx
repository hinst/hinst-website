import { useNavigate, useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { PostHeader } from './goalHeader';
import { useEffect, useState } from 'react';
import GoalPostView from './goalPostView';

export default function GoalBrowser() {
	const navigate = useNavigate();
	const params = useParams();
	const [searchParams, setSearchParams] = useSearchParams();
	const id: string = params.id!;
	const [activePostDate, setActivePostDate] = useState<string>(searchParams.get('activePostDate') || '');

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

	return <div style={{display: 'flex', gap: 20, minHeight: 0}}>
		<div style={{
			display: 'flex',
			overflowY: 'auto',
			paddingRight: 10,
			flexShrink: 0,
			flexBasis: 'fit-content',
		}}>
			<GoalCalendarPanel
				id={id}
				receivePosts={receivePosts}
				activePostDate={activePostDate}
			/>
		</div>
		<div style={{display: 'flex', overflowY: 'auto', flexGrow: 1}}>
			{ activePostDate
				? <GoalPostView
					goalId={id}
					postDate={activePostDate}
					style={{maxWidth: 1000}}
				/>
				: undefined }
		</div>
	</div>;
}
import { useNavigate, useParams, useSearchParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { PostHeader } from './goalHeader';
import { useState } from 'react';
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

	return <div style={{display: 'flex', gap: 20}}>
		<GoalCalendarPanel id={id} receivePosts={receivePosts} activePostDate={activePostDate}/>
		{ activePostDate
			? <GoalPostView
				goalId={id}
				postDate={activePostDate}
				style={{maxWidth: 1000}}
			  />
			: undefined }
	</div>;
}
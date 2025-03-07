import { useNavigate, useParams } from 'react-router';
import GoalCalendarPanel from './goalCalendarPanel';
import { PostHeader } from './goalHeader';
import { useState } from 'react';
import GoalPostView from './goalPostView';

export default function GoalBrowser() {
	const navigate = useNavigate();
	const params = useParams();
	const id: string = params.id!;
	const [activePostDate, setActivePostDate] = useState<string>('');

	function receivePosts(posts: PostHeader[]) {
		console.log(posts);
		if (posts.length && !activePostDate) {
			setActivePostDate(posts[posts.length - 1].date);
		}
	}

	return <div style={{display: 'flex'}}>
		<GoalCalendarPanel id={id} receivePosts={receivePosts}/>
		{ activePostDate ? <GoalPostView postDate={activePostDate}/> : undefined }
	</div>;
}
import { useEffect, useRef, useState } from 'react';
import { API_URL } from '../api';
import { PostHeader } from './goalHeader';
import GoalCalendar from './goalCalendar';

export default function GoalCalendarPanel(
	props: {
		id: string,
		activePostDate: string,
		receivePosts?: (posts: PostHeader[]) => void,
	}
) {
	const [isLoading, setIsLoading] = useState(0);
	const isLoadingRef = useRef(0);
	isLoadingRef.current = isLoading;

	const [posts, setPosts] = useState(new Array<PostHeader>());

	async function loadPosts() {
		setIsLoading(isLoadingRef.current + 1);
		try {
			const response = await fetch(API_URL + '/goalPosts?id=' + encodeURIComponent(props.id));
			if (!response.ok)
				throw new Error(response.statusText);
			const posts = await response.json();
			setPosts(posts);
			if (props.receivePosts)
				props.receivePosts(posts);
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(() => { loadPosts() }, []);

	return <div>
		{ isLoading
			? <div className='ms-loading'></div>
			: <GoalCalendar posts={posts} activePostDate={props.activePostDate} /> }
	</div>;
}

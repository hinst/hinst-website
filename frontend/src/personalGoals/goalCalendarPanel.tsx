import { useEffect, useRef, useState } from 'react';
import { API_URL } from 'src/typescript/global';
import { PostHeader } from './goalHeader';
import GoalCalendar from './goalCalendar';

export default function GoalCalendarPanel(
	props: {
		id: string,
		activePostDate: string,
		receivePosts?: (posts: PostHeader[]) => void,
		reload: number,
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
			let posts: PostHeader[]= await response.json();
			posts = posts.filter(post => post.type === 'post');
			setPosts(posts);
			if (props.receivePosts)
				props.receivePosts(posts);
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(
		() => { loadPosts() },
		[props.reload]
	);

	return <div>
		{ isLoading ? <div className='ms-loading'/> : undefined }
		<GoalCalendar posts={posts} activePostDate={props.activePostDate} />
	</div>;
}

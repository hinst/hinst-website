import { useEffect, useRef, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { GoalPostHeaderWithMethods } from 'src/typescript/restObjectExtensions';
import GoalCalendar from './goalCalendar';

export default function GoalCalendarPanel(props: {
	id: string;
	activePostDate: number;
	receivePosts?: (posts: GoalPostHeaderWithMethods[]) => void;
	reload: number;
}) {
	const [isLoading, setIsLoading] = useState(0);
	const isLoadingRef = useRef(0);
	isLoadingRef.current = isLoading;

	const [posts, setPosts] = useState(new Array<GoalPostHeaderWithMethods>());

	async function loadPosts() {
		setIsLoading(isLoadingRef.current + 1);
		try {
			const posts = await apiClient.getGoalPosts(parseInt(props.id) || 0);
			setPosts(posts);
			if (props.receivePosts) props.receivePosts(posts);
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(() => {
		loadPosts();
	}, [props.reload]);

	return (
		<div>
			{isLoading ? <div className='ms-loading' /> : undefined}
			<GoalCalendar posts={posts} activePostDate={props.activePostDate} />
		</div>
	);
}

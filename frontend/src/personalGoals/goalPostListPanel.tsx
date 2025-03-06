import { useEffect, useRef, useState } from 'react';
import { API_URL } from '../api';
import { useParams } from 'react-router';
import { PostHeader } from './goalHeader';

export default function GoalPostListPanel() {
	const [isLoading, setIsLoading] = useState(0);
	const isLoadingRef = useRef(0);
	isLoadingRef.current = isLoading;

	const params = useParams();
	const id: string = params.id!;

	const [posts, setPosts] = useState(new Array<PostHeader>());

	async function loadPosts() {
		setIsLoading(isLoadingRef.current + 1);
		try {
			const response = await fetch(API_URL + '/goalPosts?id=' + encodeURIComponent(id));
			if (!response.ok)
				throw new Error(response.statusText);
			setPosts(await response.json());
		} finally {
			setIsLoading(isLoadingRef.current - 1);
		}
	}
	useEffect(() => { loadPosts() }, []);

	function getMonthlyPosts() {
		
	}

	return <div>
		{ isLoading ? <div className='ms-loading'></div> : undefined }
		{posts.map(post =>
			<div key={post.id}>
				{post.date}
			</div>
		)}
	</div>;
}
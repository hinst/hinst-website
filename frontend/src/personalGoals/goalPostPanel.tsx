import { useEffect, useState } from 'react';
import { SmartPostExtended, SmartPostImage } from './smartPost';
import { API_URL } from '../global';
import GoalPostView from './goalPostView';
import { Tool } from 'react-feather';

export default function GoalPostPanel(props: {
	goalId: string,
	postDate: string,
	goalManagerMode: boolean,
}) {
	const [isLoading, setIsLoading] = useState(false);
	const [postData, setPostData] = useState<SmartPostExtended | undefined>(undefined);

	async function load() {
		setIsLoading(true);
		setPostData(undefined);
		try {
			const response = await fetch(API_URL + '/goalPost' +
				'?goalId=' + encodeURIComponent(props.goalId) +
				'&postDateTime=' + props.postDate);
			const postData: SmartPostExtended = await response.json();
			setPostData(postData);
			loadImages();
		} finally {
			setIsLoading(false);
		}
	};

	async function loadImages() {
		const response = await fetch(API_URL + '/goalPost/images' +
			'?goalId=' + encodeURIComponent(props.goalId) +
			'&postDateTime=' + props.postDate);
		const images: SmartPostImage[] = await response.json();
		setPostData(postData => postData ? {...postData, images} : undefined);
	}

	useEffect(() => {
		load();
	}, [props.goalId, props.postDate]);

	async function setPublic(isPublic: boolean) {
		const url = API_URL + '/goalPost/setPublic' +
			'?goalId=' + encodeURIComponent(props.goalId) +
			'&postDateTime=' + encodeURIComponent(props.postDate) +
			'&isPublic=' + encodeURIComponent('' + isPublic);
		const response = await fetch(url);
		if (!response.ok)
			throw new Error('Cannot update post visibility. Status: ' + response.statusText);
		if (postData)
			setPostData({...postData, isPublic});
	}

	function RenderGoalManagement() {
		return <div
			className='ms-alert ms-light'
			style={{
				display: 'flex',
				gap: 10,
				alignItems: 'center',
			}}
		>
			<Tool/>
			<input
				type='checkbox'
				checked={postData?.isPublic}
				onChange={() => setPublic(!postData?.isPublic)}
			/>
			public
		</div>;
	}

	return <div>
		{ isLoading ? <div className='ms-loading' /> : undefined }
		{ postData ? [
			RenderGoalManagement(),
			<GoalPostView postData={postData} />
		] : undefined }
	</div>;
}

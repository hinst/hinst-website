import { useEffect, useState } from 'react';
import { SmartPostExtended, SmartPostImage } from './smartPost';
import { API_URL } from '../global';
import GoalPostView from './goalPostView';

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

	return <div>
		{ isLoading ? <div className='ms-loading' /> : undefined }
		{ postData ? <GoalPostView postData={postData} />: undefined }
	</div>;
}

function GoalImage(props: { data: string }) {
	return <a href={props.data}>
		<img
			className='ms-card ms-border'
			width={240}
			height={240}
			style={{
				width: 240,
				height: 240,
				objectFit: 'cover',
				margin: 0,
				padding: 0,
			}}
			src={props.data}
			alt='Image'
		/>
	</a>;
}
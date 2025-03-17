import { useEffect, useState } from 'react';
import { SmartPostExtended, SmartPostImage } from './smartPost';
import { API_URL } from '../global';
import SafeHtmlView from '../safeHtmlView';
import { Info } from 'react-feather';

export default function GoalPostView(props: {
	goalId: string,
	postDate: string,
	style?: React.CSSProperties
}) {
	const [isLoading, setIsLoading] = useState(false);
	const [postData, setPostData] = useState<SmartPostExtended | undefined>(undefined);
	const [images, setImages] = useState<string[]>([]);

	async function load() {
		setIsLoading(true);
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
		setImages(images?.map(image => image.dataUrl) || []);
	}

	useEffect(() => {
		load();
	}, [props.goalId, props.postDate]);

	return <div style={props.style}>
		{ isLoading ? <div className='ms-loading' /> : undefined }
		{postData
			? <div>
				{postData.isAutoTranslated
					? <div className='ms-alert ms-light' style={{display: 'flex', alignItems: 'center'}}>
						<Info/> &nbsp; This text was automatically converted to {postData.languageName} language using AI translator tool.
					  </div>
					: undefined
				}
				{postData.languageNamePending
					? <div className='ms-alert ms-light' style={{display: 'flex', alignItems: 'center'}}>
						<Info className='ms-text-secondary'/>
						&nbsp;
						The automatic translation of this text to {postData.languageNamePending} language is not available yet.
						Please come back later or check older posts.
					  </div>
					: undefined
				}
				<SafeHtmlView htmlText={postData.msg} />
			</div>
			: undefined}
		<div
			style={{
				display: 'flex',
				flexWrap: 'wrap',
				gap: 10,
			}}
		>
			{images.map(image => <GoalImage key={image} data={image} />)}
		</div>
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
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
		setImages(images.map(image => image.dataUrl));
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
						<Info/> &nbsp; This text was automatically converted into {postData.languageName} language using AI translator tool.
					  </div>
					: undefined
				}
				<SafeHtmlView htmlText={postData.msg} />
				{images.map((image, index) => <div key={index} style={{marginTop: '10px'}}>
					<img src={image} alt='' style={{maxWidth: '100%'}} />
				</div>)}
			</div>
			: undefined}
	</div>;
}
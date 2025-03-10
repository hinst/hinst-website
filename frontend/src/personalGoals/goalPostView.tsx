import { useEffect, useRef, useState } from 'react';
import { SmartPost } from './smartPost';
import { API_URL } from '../global';
import SafeHtmlView from '../safeHtmlView';

export default function GoalPostView(props: {
	goalId: string,
	postDate: string,
	style?: React.CSSProperties
}) {
	const [isLoading, setIsLoading] = useState(false);
	const [postData, setPostData] = useState<SmartPost | undefined>(undefined);

	async function load() {
		setIsLoading(true);
		try {
			const response = await fetch(API_URL + '/goalPost' +
				'?goalId=' + encodeURIComponent(props.goalId) +
				'&postDateTime=' + props.postDate);
			const postData = await response.json();
			setPostData(postData);
		} finally {
			setIsLoading(false);
		}
	};

	useEffect(() => {
		load();
	}, [props.goalId, props.postDate]);

	return <div style={props.style}>
		{ isLoading ? <div className='ms-loading' /> : undefined }
		{postData
			? <SafeHtmlView htmlText={postData.msg} />
			: undefined}
	</div>;
}
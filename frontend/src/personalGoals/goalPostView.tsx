import { useEffect, useState } from 'react';
import { SmartPost } from './smartPost';
import { API_URL } from '../api';
import SafeHtmlView from '../safeHtmlView';

export default function GoalPostView(props: {
	goalId: string,
	postDate: string,
	style?: React.CSSProperties
}) {
	const [postData, setPostData] = useState<SmartPost | undefined>(undefined);

	async function load() {
		const response = await fetch(API_URL + '/goalPost' +
			'?goalId=' + encodeURIComponent(props.goalId) +
			'&postDateTime=' + props.postDate);
		const postData = await response.json();
		setPostData(postData);
	};

	useEffect(() => {
		load();
	}, []);

	return <div style={props.style}>
		{postData
			? <SafeHtmlView htmlText={postData.msg} />
			: undefined}
	</div>;
}
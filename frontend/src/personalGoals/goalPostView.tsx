import { useEffect, useState } from 'react';
import { SmartPost } from './smartPost';
import { API_URL } from '../api';

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
			? <div dangerouslySetInnerHTML={{__html: postData.msg}} />
			: undefined}
	</div>;
}
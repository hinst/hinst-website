import { useEffect, useState } from 'react';
import { SmartPost } from './smartPost';
import { API_URL } from '../api';

export default function GoalPostView(props: {
	goalId: string,
	postDate: string,
	style?: React.CSSProperties
}) {
	const [postData, setPostData] = useState<SmartPost | undefined>(undefined);
	const [html, setHtml] = useState<string | undefined>(undefined);

	async function load() {
		const response = await fetch(API_URL + '/goalPost' +
			'?goalId=' + encodeURIComponent(props.goalId) +
			'&postDateTime=' + props.postDate);
		const postData = await response.json();
		setPostData(postData);

		const parser = new DOMParser();
		const parsedHTML = parser.parseFromString(postData.msg, 'text/html');
		const scriptTags = parsedHTML.getElementsByTagName('script');
		[...scriptTags].forEach(tag => tag.remove());
		setHtml(parsedHTML.body.innerHTML);
	};

	useEffect(() => {
		load();
	}, []);

	return <div style={props.style}>
		{html
			? <div dangerouslySetInnerHTML={{__html: html}} />
			: undefined}
	</div>;
}
import { useEffect, useState } from 'react';
import { SmartPostExtended } from './smartPost';
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

	async function load() {
		setIsLoading(true);
		try {
			const response = await fetch(API_URL + '/goalPost' +
				'?goalId=' + encodeURIComponent(props.goalId) +
				'&postDateTime=' + props.postDate);
			const postData: SmartPostExtended = await response.json();
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
			? <div>
				{postData.isAutoTranslated
					? <div className='ms-alert ms-light' style={{display: 'flex', alignItems: 'center'}}>
						<Info/> &nbsp; This text was automatically converted into {postData.languageName} language using AI translator tool.
					  </div>
					: undefined
				}
				<SafeHtmlView htmlText={postData.msg} />
			</div>
			: undefined}
	</div>;
}
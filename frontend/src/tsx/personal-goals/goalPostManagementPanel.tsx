import { Tool } from 'react-feather';
import { GoalPostObject, GoalPostObjectExtended } from 'src/typescript/personal-goals/smartPost';
import { API_URL } from 'src/typescript/global';
import { useState } from 'react';

export default function GoalPostManagementPanel(props: {
	postData: GoalPostObjectExtended;
	setPostData: (postData: GoalPostObjectExtended) => void;
	onChange: () => void;
}) {
	const [isLoading, setIsLoading] = useState(false);
	async function setPublic(isPublic: boolean) {
		if (isLoading) return;
		setIsLoading(true);
		try {
			const url =
				API_URL +
				'/goalPost/setPublic' +
				'?goalId=' +
				encodeURIComponent(props.postData.goalId) +
				'&postDateTime=' +
				encodeURIComponent(props.postData.dateTime) +
				'&isPublic=' +
				encodeURIComponent('' + isPublic);
			const response = await fetch(url);
			if (!response.ok)
				throw new Error('Cannot update post visibility. Status: ' + response.statusText);
			props.setPostData({ ...props.postData, isPublic });
			props.onChange();
		} finally {
			setIsLoading(false);
		}
	}
	return (
		<div
			className='ms-alert ms-light'
			style={{
				display: 'flex',
				gap: 10,
				alignItems: 'center'
			}}
		>
			<Tool />
			<input
				disabled={isLoading}
				type='checkbox'
				checked={props.postData?.isPublic}
				onChange={() => setPublic(!props.postData?.isPublic)}
			/>
			public
		</div>
	);
}

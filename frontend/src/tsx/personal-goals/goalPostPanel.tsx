import { useContext, useEffect, useState } from 'react';
import { GoalPostObject, SmartPostImage } from 'src/typescript/personal-goals/smartPost';
import { API_URL } from 'src/typescript/global';
import GoalPostView from './goalPostView';
import GoalPostManagementPanel from './goalPostManagementPanel';
import { AppContext } from 'src/tsx/context';

export default function GoalPostPanel(props: {
	goalId: string;
	postDate: string;
	goalManagerMode: boolean;
	onChange: () => void;
}) {
	const [isLoading, setIsLoading] = useState(false);
	const [postData, setPostData] = useState<GoalPostObject | undefined>(undefined);
	const [errorMessage, setErrorMessage] = useState<string>('');
	const context = useContext(AppContext);

	async function load() {
		setIsLoading(true);
		setPostData(undefined);
		try {
			const response = await fetch(
				API_URL +
					'/goalPost' +
					'?goalId=' +
					encodeURIComponent(props.goalId) +
					'&postDateTime=' +
					props.postDate
			);
			if (!response.ok) {
				setErrorMessage(response.statusText);
				return;
			}
			const postData: GoalPostObject = await response.json();
			setPostData(postData);
			loadImages();
		} finally {
			setIsLoading(false);
		}
	}

	async function loadImages() {
		const response = await fetch(
			API_URL +
				'/goalPost/images' +
				'?goalId=' +
				encodeURIComponent(props.goalId) +
				'&postDateTime=' +
				props.postDate
		);
		const images: SmartPostImage[] = await response.json();
		setPostData((postData) => (postData ? { ...postData, images } : undefined));
	}

	useEffect(() => {
		load();
	}, [props.goalId, props.postDate]);

	return (
		<div>
			{isLoading ? <div className='ms-loading' /> : undefined}
			{postData ? (
				<>
					{context.goalManagerMode ? (
						<GoalPostManagementPanel
							postData={postData}
							setPostData={setPostData}
							onChange={props.onChange}
						/>
					) : undefined}
					<GoalPostView postData={postData} />
				</>
			) : (
				errorMessage
			)}
		</div>
	);
}

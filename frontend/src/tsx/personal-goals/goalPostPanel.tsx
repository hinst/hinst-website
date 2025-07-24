import { useContext, useEffect, useState } from 'react';
import { GoalPostObject } from 'src/typescript/personal-goals/smartPost';
import GoalPostView from './goalPostView';
import GoalPostManagementPanel from './goalPostManagementPanel';
import { AppContext } from 'src/tsx/context';
import { apiClient } from 'src/typescript/apiClient';

export default function GoalPostPanel(props: {
	goalId: number;
	postDate: number;
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
			let postData: GoalPostObject;
			try {
				postData = await apiClient.getGoalPost(props.goalId, props.postDate);
			} catch (e) {
				setErrorMessage((e as Error).message);
				return;
			}
			setPostData(postData);
		} finally {
			setIsLoading(false);
		}
	}

	function receiveChange() {
		load();
		props.onChange();
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
							onChange={receiveChange}
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

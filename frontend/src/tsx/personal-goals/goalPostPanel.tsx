import { useContext, useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import type { GoalPostObject } from 'src/typescript/apiTypes';
import { AppContext } from 'src/typescript/appContext';
import GoalPostManagementPanel from './goalPostManagementPanel';
import GoalPostView from './goalPostView';

export default function GoalPostPanel(props: {
	goalId: number;
	postDate: number;
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
					{context.isAdminMode ? (
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

import { Edit, Tool } from 'react-feather';
import { GoalPostObject, GoalPostObjectExtended } from 'src/typescript/personal-goals/smartPost';
import { useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';

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
			const response = await apiClient.goalPostSetPublic(
				props.postData.goalId,
				props.postData.dateTime,
				isPublic
			);
			if (!response.ok)
				throw new Error('Cannot update post visibility. Status: ' + response.statusText);
			props.setPostData({ ...props.postData, isPublic });
			props.onChange();
		} finally {
			setIsLoading(false);
		}
	}

	const [content, setContent] = useState<GoalPostObject | undefined>(undefined);

	async function loadContent() {
		if (isLoading) return;
		setIsLoading(true);
		try {
			setContent(await apiClient.getGoalPost(props.postData.goalId, props.postData.dateTime));
		} finally {
			setIsLoading(false);
		}
	}

	async function saveContent() {
		
	}

	function toggleEditMode() {
		if (content) setContent(undefined);
		else loadContent();
	}

	return (
		<div
			style={{ display: 'flex', flexDirection: 'column', gap: 10 }}
			className='ms-alert ms-light'
		>
			<div
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
			<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
				<div>
					<button
						type='button'
						className='ms-btn ms-action'
						disabled={isLoading}
						style={{ display: 'flex', gap: 10, alignItems: 'center' }}
						onClick={toggleEditMode}
					>
						<Edit /> Edit content
					</button>
				</div>
				{content ? (
					<textarea
						style={{ fontFamily: 'monospace' }}
						rows={20}
						value={content?.text}
						onChange={(event) => setContent({ ...content, text: event.target.value })}
					/>
				) : undefined}
			</div>
		</div>
	);
}

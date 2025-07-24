import { Edit, Save, Tool } from 'react-feather';
import { GoalPostObject, GoalPostObjectExtended } from 'src/typescript/personal-goals/smartPost';
import { useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';

export default function GoalPostManagementPanel(props: {
	postData: GoalPostObject;
	setPostData: (postData: GoalPostObject) => void;
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

	async function saveText() {
		if (isLoading || !content) return;
		setIsLoading(true);
		try {
			const response = await apiClient.setGoalPostText(
				props.postData.goalId,
				props.postData.dateTime,
				content.languageTag,
				content.text
			);
			if (!response.ok) throw new Error('Cannot save text. Status: ' + response.statusText);
			props.onChange();
		} finally {
			setIsLoading(false);
		}
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
				<div style={{ display: 'flex', gap: 10 }}>
					<button
						type='button'
						className='ms-btn ms-action'
						disabled={isLoading}
						style={{ display: 'flex', gap: 10, alignItems: 'center' }}
						onClick={toggleEditMode}
					>
						<Edit /> Edit content
					</button>
					{content ? (
						<button
							type='button'
							className='ms-btn ms-action2'
							disabled={isLoading}
							style={{ display: 'flex', gap: 10, alignItems: 'center' }}
							onClick={saveText}
						>
							<Save /> Save text
						</button>
					) : undefined}
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

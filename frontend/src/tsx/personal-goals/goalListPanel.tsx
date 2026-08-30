import { useEffect, useState } from 'react';
import { GoalObjectEx } from 'src/typescript/rest_objects/goalObjectEx';
import GoalList from './goalList';
import { apiClient } from 'src/typescript/apiClient';

export default function GoalListPanel() {
	const [isLoading, setIsLoading] = useState(false);
	const [goals, setGoals] = useState(new Array<GoalObjectEx>());
	async function loadGoals() {
		setIsLoading(true);
		try {
			setGoals(await apiClient.getGoals());
		} finally {
			setIsLoading(false);
		}
	}
	useEffect(() => {
		loadGoals();
	}, []);
	return (
		<div>
			{isLoading ? <div className='ms-loading'></div> : undefined}
			<GoalList goals={goals} />
		</div>
	);
}

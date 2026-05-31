import { useEffect, useState } from 'react';
import { GoalHeader } from 'src/typescript/personal-goals/goalRecord';
import { API_URL } from 'src/typescript/global';
import GoalList from './goalList';
import { apiClient } from 'src/typescript/apiClient';

export default function GoalListPanel() {
	const [isLoading, setIsLoading] = useState(false);
	const [goals, setGoals] = useState(new Array<GoalHeader>());
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

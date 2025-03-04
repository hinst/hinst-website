import { useEffect, useState } from 'react';
import { GoalHeader } from './goalHeader';
import { API_URL } from '../api';
import GoalList from './goalList';

export default function GoalListPanel() {
	const [isLoading, setIsLoading] = useState(false);
	const [goals, setGoals] = useState(new Array<GoalHeader>());
	async function loadGoals() {
		setIsLoading(true);
		try {
			const response = await fetch(API_URL + '/goals');
			if (!response.ok)
				throw new Error(response.statusText);
			const data = await response.json();
			setGoals(data);
		} finally {
			setIsLoading(false);
		}
	};
	useEffect(
		() => { loadGoals(); },
		[],
	);
	return <div
		className='ms-card ms-border'
		style={{display: 'flex', flexDirection: 'column', width: 'fit-content', maxWidth: '100%'}}
	>
		<div className='ms-card-title'>
			My personal goals
		</div>
		{ isLoading ? <div className='ms-loading'></div> : undefined }
		<GoalList goals={goals}/>
	</div>;
}
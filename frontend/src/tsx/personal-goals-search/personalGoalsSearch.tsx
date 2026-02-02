import { useState } from 'react';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { SearchBar } from './searchBar';
import { ItemRow } from './itemRow';
import { apiClient } from 'src/typescript/apiClient';

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostRecord>>([]);
	async function search(queryText: string) {
		const items = await apiClient.searchGoalPosts(queryText);
		setItems(items);
	}
	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<SearchBar onSearch={search} />
			<div>Results: {items.length}</div>
			<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
				{items.map((item, index) => (
					<ItemRow key={index} item={item} />
				))}
			</div>
		</div>
	);
}

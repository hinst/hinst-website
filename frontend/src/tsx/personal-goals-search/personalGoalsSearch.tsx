import { useState } from 'react';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { SearchBar } from './searchBar';
import { ItemRow } from './itemRow';

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostRecord>>([
		new GoalPostRecord(0, 0, true, 'type', 'hello')
	]);
	function search(queryText: string) {
		console.log('Searching for:', queryText);
	}
	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<SearchBar onSearch={search} />
			<div>Results: {items.length}</div>
			<div>
				{items.map((item, index) => (
					<ItemRow key={index} item={item} />
				))}
			</div>
		</div>
	);
}

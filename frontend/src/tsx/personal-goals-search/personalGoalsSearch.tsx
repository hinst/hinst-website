import { useEffect, useState } from 'react';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { SearchBar } from './searchBar';
import { ItemRow } from './itemRow';
import { apiClient } from 'src/typescript/apiClient';
import { useSearchParams } from 'react-router';

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostRecord>>([]);
	const [searchParams, setSearchParams] = useSearchParams();

	function getQuery() {
		return searchParams.get('query');
	}

	function goSearch(query: string) {
		setSearchParams({ query: query }, { replace: true });
	}

	async function search(query: string) {
		const items = await apiClient.searchGoalPosts(query);
		setItems(items);
	}

	useEffect(() => {
		const query = getQuery();
		if (query) search(query);
	}, [getQuery()]);

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<SearchBar onSearch={goSearch} text={getQuery() || ''} />
			<div>Results: {items.length}</div>
			<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
				{items.map((item, index) => (
					<ItemRow key={index} item={item} />
				))}
			</div>
		</div>
	);
}

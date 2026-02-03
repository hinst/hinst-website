import { useEffect, useState } from 'react';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { SearchBar } from './searchBar';
import { ItemRow } from './itemRow';
import { apiClient } from 'src/typescript/apiClient';
import { useSearchParams } from 'react-router';

const SERVER_SIDE_LIMIT = 100;

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostRecord>>([]);
	const [searchParams, setSearchParams] = useSearchParams();
	const [isLoading, setIsLoading] = useState(false);

	function getQuery() {
		return searchParams.get('query');
	}

	function goSearch(query: string) {
		setSearchParams({ query: query }, { replace: true });
	}

	async function search(query: string) {
		if (isLoading) return;
		setIsLoading(true);
		try {
			const items = await apiClient.searchGoalPosts(query);
			setItems(items);
		} finally {
			setIsLoading(false);
		}
	}

	useEffect(() => {
		const query = getQuery();
		if (query) search(query);
	}, [getQuery()]);

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<div style={{ display: 'flex', gap: 5, alignItems: 'center' }}>
				{isLoading ? <div className='ms-loading' /> : undefined}
				<SearchBar onSearch={goSearch} text={getQuery() || ''} disabled={isLoading} />
			</div>
			<div>
				Results: {items.length}
				{items.length === SERVER_SIDE_LIMIT ? '...' : ''}
			</div>
			<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
				{items.map((item, index) => (
					<ItemRow key={index} item={item} />
				))}
			</div>
		</div>
	);
}

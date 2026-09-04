import { useContext, useEffect, useState } from 'react';
import { useSearchParams } from 'react-router';
import { apiClient } from 'src/typescript/apiClient';
import { AppContext } from 'src/typescript/appContext';
import { PageTitle } from 'src/typescript/pageTitle';
import type { GoalPostHeaderEx } from 'src/typescript/rest_objects/goalPostHeaderEx';
import { ItemRow } from './itemRow';
import { SearchBar } from './searchBar';

const SERVER_SIDE_LIMIT = 100;

export function PersonalGoalsSearch() {
	const context = useContext(AppContext);
	const [items, setItems] = useState<Array<GoalPostHeaderEx>>([]);
	const [searchParams, setSearchParams] = useSearchParams();
	const [isLoading, setIsLoading] = useState(false);

	function getQuery() {
		return searchParams.get('query');
	}

	function goSearch(query: string) {
		setSearchParams({ query: query }, { replace: true });
	}

	function updateTitle(items: unknown[]) {
		const secondaryTitle =
			'Results: ' + items.length + (items.length === SERVER_SIDE_LIMIT ? '+' : '');
		context.setPageTitle(new PageTitle('Search', secondaryTitle));
	}

	async function search(query: string) {
		if (isLoading) return;
		setIsLoading(true);
		try {
			const items = await apiClient.searchGoalPosts(query);
			setItems(items || []);
			updateTitle(items);
		} finally {
			setIsLoading(false);
		}
	}

	useEffect(() => {
		updateTitle(items);
	}, []);

	useEffect(() => {
		const query = getQuery();
		if (query) {
			const _ignoredPromise = search(query);
		}
	}, [getQuery()]);

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<div style={{ display: 'flex', gap: 5, alignItems: 'center' }}>
				{isLoading ? <div className='ms-loading' /> : undefined}
				<SearchBar onSearch={goSearch} text={getQuery() || ''} disabled={isLoading} />
			</div>
			<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
				{items.map((item, index) => (
					<ItemRow key={index} item={item} />
				))}
			</div>
		</div>
	);
}

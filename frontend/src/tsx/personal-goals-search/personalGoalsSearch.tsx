import { useState } from 'react';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';
import { SearchBar } from './searchBar';

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
					<div key={index}>
						<a
							href={`/personal-goals/${item.goalId}/${item.dateTime}.html`}
							style={{ display: 'inline-flex', gap: 10 }}
						>
							<button
								type='button'
								className='ms-btn ms-primary ms-outline'
								style={{
									fontFamily: 'monospace',
									minWidth: 50,
									padding: 8,
									display: 'flex',
									justifyContent: 'center'
								}}
							>
								<span>{item.dateTime}</span>
							</button>
							<div
								className='ms-text-main'
								style={{ display: 'flex', alignItems: 'center' }}
							>
								{item.title}
							</div>
						</a>
					</div>
				))}
			</div>
		</div>
	);
}

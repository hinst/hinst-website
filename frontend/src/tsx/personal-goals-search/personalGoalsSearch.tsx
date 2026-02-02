import { useState } from 'react';
import { Search } from 'react-feather';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostRecord>>([
		new GoalPostRecord(0, 0, true, 'type', 'hello')
	]);
	const [queryText, setQueryText] = useState<string>('');
	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<div style={{ display: 'flex', gap: 5 }}>
				<input
					type='text'
					placeholder='Search text...'
					autoFocus={true}
					value={queryText}
					onChange={(e) => setQueryText(e.target.value)}
				/>
				<button type='button' style={{ display: 'flex', gap: 5, alignItems: 'center' }}>
					<Search />
					Search
				</button>
			</div>
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

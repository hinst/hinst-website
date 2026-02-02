import { Search } from 'react-feather';

export function PersonalGoalsSearch() {
	return (
		<div>
			<div style={{ display: 'flex', gap: 5 }}>
				<input type='text' placeholder='Search text...' autoFocus={true} />
				<button type='button' style={{ display: 'flex', gap: 5, alignItems: 'center' }}>
					<Search />
					Search
				</button>
			</div>
		</div>
	);
}

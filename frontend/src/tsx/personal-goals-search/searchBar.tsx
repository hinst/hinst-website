import { Search } from 'react-feather';
import { useState } from 'react';

interface SearchBarProps {
	onSearch: (queryText: string) => void;
}

export function SearchBar({ onSearch }: SearchBarProps) {
	const [queryText, setQueryText] = useState('');
	return (
		<div style={{ display: 'flex', gap: 5 }}>
			<input
				type='text'
				placeholder='Search text...'
				autoFocus={true}
				value={queryText}
				onChange={(e) => setQueryText(e.target.value)}
				onKeyDown={(e) => {
					if (e.key === 'Enter') onSearch(queryText);
				}}
			/>
			<button
				type='button'
				style={{ display: 'flex', gap: 5, alignItems: 'center' }}
				onClick={() => onSearch(queryText)}
			>
				<Search />
				Search
			</button>
		</div>
	);
}

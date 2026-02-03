import { Search } from 'react-feather';
import { useState } from 'react';

interface Props {
	text: string;
	disabled: boolean;
	onSearch: (queryText: string) => void;
}

export function SearchBar(props: Props) {
	const [queryText, setQueryText] = useState(props.text || '');
	return (
		<div style={{ display: 'flex', gap: 5, flexGrow: 1 }}>
			<input
				type='text'
				placeholder='Search text...'
				autoFocus={true}
				value={queryText}
				onChange={(e) => setQueryText(e.target.value)}
				onKeyDown={(e) => {
					if (e.key === 'Enter') props.onSearch(queryText);
				}}
			/>
			<button
				type='button'
				style={{ display: 'flex', gap: 5, alignItems: 'center' }}
				onClick={() => props.onSearch(queryText)}
				disabled={props.disabled}
			>
				<Search />
				Search
			</button>
		</div>
	);
}

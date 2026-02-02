import { useState } from 'react';
import { Search } from 'react-feather';
import { GoalPostObject } from 'src/typescript/personal-goals/smartPost';

export function PersonalGoalsSearch() {
	const [items, setItems] = useState<Array<GoalPostObject>>([]);
	const [queryText, setQueryText] = useState<string>('');
	return (
		<div>
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
			<div>
				<div
					style={{
						width: 'fit-content',
						padding: 0,
						margin: 0,
						overflow: 'hidden'
					}}
					className='ms-card ms-border grayscale'
				>
					<a
						href='{{$.WebPath}}/personal-goals/{{.Id}}{{$.HtmlExtension}}'
						className='ms-text-main'
					>
						<div
							style={{ display: 'flex', flexDirection: 'column' }}
							className='ms-bg-light'
						>
							<img
								width='200'
								height='100'
								src='{{$.StaticPath}}{{.Image}}'
								alt='{{.Title}}'
							/>
							<div style={{ padding: 8, maxWidth: 200 }}>Title</div>
						</div>
					</a>
				</div>
			</div>
		</div>
	);
}

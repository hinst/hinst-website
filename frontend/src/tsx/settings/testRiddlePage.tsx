import { useEffect } from 'react';

export default function TestRiddlePage(props: { setPageTitle: (title: string) => void }) {
	useEffect(() => {
		props.setPageTitle('Test Riddle');
	}, []);

	

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
			<div>
				<button
					className='ms-btn'
					onClick={() => {
						window.location.href = '/riddle/test';
					}}
				>
					TEST
				</button>
			</div>
		</div>
	);
}
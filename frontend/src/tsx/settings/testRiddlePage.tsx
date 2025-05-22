import { useEffect } from 'react';
import { apiClient } from 'src/typescript/apiClient';

export default function TestRiddlePage(props: { setPageTitle: (title: string) => void }) {
	useEffect(() => {
		props.setPageTitle('Test Riddle');
	}, []);

	async function test() {
		const riddle = await apiClient.createRiddle();
		const primeNumbers = await apiClient.getPrimeNumbers();
		const answer = riddle.solve(primeNumbers);
		console.log(riddle, answer);
	}

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
			<div>
				<button className='ms-btn ms-outline ms-primary' onClick={test}>
					TEST
				</button>
			</div>
		</div>
	);
}

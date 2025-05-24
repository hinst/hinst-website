import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { RiddleItem, RiddleSolver } from 'src/typescript/riddle';

export default function TestRiddlePage(props: { setPageTitle: (title: string) => void }) {
	useEffect(() => {
		props.setPageTitle('Test Riddle');
	}, []);
	const [isLoading, setIsLoading] = useState(false);
	const [primeNumbers, setPrimeNumbers] = useState(new Array<number>());
	const [riddle, setRiddle] = useState(new RiddleItem());
	const [haveAnswer, setHaveAnswer] = useState(false);
	const [answer, setAnswer] = useState(new Array<number>());
	const [answerCallCount, setAnswerCallCount] = useState(0);
	const [answerCount, setAnswerCount] = useState(0);

	function getProduct() {
		return answer.reduce((product, item) => (product * item) % riddle.limit, 1);
	}

	async function test() {
		setIsLoading(true);
		try {
			const riddle = await apiClient.createRiddle();
			setRiddle(riddle);
			const primeNumbers = await apiClient.getPrimeNumbers();
			setPrimeNumbers(primeNumbers);

			const solver = new RiddleSolver(
				primeNumbers,
				riddle.product,
				riddle.steps,
				riddle.limit
			);
			setHaveAnswer(solver.solve());
			setAnswer(solver.sequence);
			setAnswerCallCount(solver.callCount);

			solver.callCount = 0;
			setAnswerCount(solver.count());
		} finally {
			setIsLoading(false);
		}
	}

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
			<div style={{ display: 'flex', flexDirection: 'row', gap: 10 }}>
				<button
					className='ms-btn ms-outline ms-primary'
					onClick={test}
					disabled={isLoading}
				>
					TEST
				</button>
				{isLoading ? <div className='ms-loading' /> : undefined}
			</div>
			<ul style={{ margin: 0 }}>
				<li>
					riddle: <pre>{JSON.stringify(riddle)}</pre>
				</li>
				<li>
					prime numbers: <pre>{primeNumbers.length}</pre>
				</li>
				<li>
					answer: <b>{haveAnswer ? 'true' : 'false'}</b> {getProduct()}{' '}
					{getProduct() === riddle.product ? 'ok' : 'error'}
					<pre>{answer.join(',')}</pre>
					call count: {answerCallCount}
				</li>
				<li>answer count: {answerCount}</li>
			</ul>
		</div>
	);
}

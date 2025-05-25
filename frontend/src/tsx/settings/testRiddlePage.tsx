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
	const [timeToAnswer, setTimeToAnswer] = useState(0);
	const [timeToCount, setTimeToFindAll] = useState(0);

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
			let time = new Date().getTime();
			setHaveAnswer(solver.solve());
			time = new Date().getTime() - time;
			setTimeToAnswer(time);
			setAnswer(solver.sequence);
			setAnswerCallCount(solver.callCount);

			solver.callCount = 0;
			time = new Date().getTime();
			setAnswerCount(solver.count());
			time = new Date().getTime() - time;
			setTimeToFindAll(time);
		} finally {
			setIsLoading(false);
		}
	}

	const listItemStyle: React.CSSProperties = { marginBottom: 10 };
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
				<li style={listItemStyle}>
					riddle: <code>{JSON.stringify(riddle)}</code>
				</li>
				<li style={listItemStyle}>
					prime numbers: <code>{primeNumbers.length}</code>
				</li>
				<li style={listItemStyle}>
					answer: <b>{haveAnswer ? 'true' : 'false'}</b> <code>{getProduct()}</code>
					{getProduct() === riddle.product ? 'ok' : 'error'}
					<pre>{answer.join(',')}</pre>
					call count: {answerCallCount}
				</li>
				<li style={listItemStyle}>
					time to answer: <code>{(timeToAnswer / 1000).toFixed(3)}</code> seconds
				</li>
				<li style={listItemStyle}>
					answer count: <code>{answerCount}</code>
				</li>
				<li style={listItemStyle}>
					time to count all answers: <code>{(timeToCount / 1000).toFixed(3)}</code>{' '}
					seconds
				</li>
			</ul>
		</div>
	);
}

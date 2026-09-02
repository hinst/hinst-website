import { useContext, useEffect } from 'react';
import GoalListPanel from './personal-goals/goalListPanel';
import { AppContext } from 'src/tsx/context';

export default function HomePage() {
	const context = useContext(AppContext);
	useEffect(() => {
		context.setPageTitle('Category list');
	}, []);
	return (
		<div>
			<GoalListPanel />
		</div>
	);
}

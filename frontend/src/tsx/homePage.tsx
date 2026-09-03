import { useContext, useEffect } from 'react';
import { AppContext } from 'src/tsx/appContext';
import GoalListPanel from './personal-goals/goalListPanel';

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

import { useEffect } from 'react';
import GoalListPanel from './personalGoals/goalListPanel';

export default function HomePage(props: {
	setPageTitle: (title: string) => void
}) {
	useEffect(() => {
		props.setPageTitle('My Personal Goals');
	}, []);
	return <div>
		<GoalListPanel />
	</div>;
}
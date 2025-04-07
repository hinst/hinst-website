import { useEffect } from 'react';
import GoalListPanel from './personal-goals/goalListPanel';

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
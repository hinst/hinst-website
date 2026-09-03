import { useContext, useEffect } from 'react';
import { AppContext } from 'src/typescript/appContext';
import { PageTitle } from 'src/typescript/pageTitle';
import GoalListPanel from './personal-goals/goalListPanel';

export default function HomePage() {
	const context = useContext(AppContext);
	useEffect(() => {
		context.setPageTitle(new PageTitle('', 'Journals'));
	}, []);
	return (
		<div>
			<GoalListPanel />
		</div>
	);
}

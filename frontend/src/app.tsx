import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalPostListPanel from './personalGoals/goalPostListPanel';

const PAGE_TITLE = 'Showcase Website';

export default function App() {
	const [pageTitle, setPageTitle] = useState(PAGE_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);
	useEffect(() => {
		if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
			document.getElementsByTagName('html')[0].setAttribute('data-theme', 'dark');
	}, []);
	return <div style={{marginLeft: 10, marginTop: 10, marginRight: 10}}>
		<HashRouter>
			<div style={{ marginBottom: 10 }}>
				<Header title={pageTitle} />
			</div>
			<Routes>
				<Route path='/' element={<HomePage />} />
				<Route path='/personal-goals/:id' element={<GoalPostListPanel />} />
			</Routes>
		</HashRouter>
	</div>;
}

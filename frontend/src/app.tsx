import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';

const PAGE_TITLE = 'Showcase Website';

export default function App() {
	const [pageTitle, setPageTitle] = useState(PAGE_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);
	return <HashRouter>
		<div style={{ marginBottom: 10 }}>
			<Header title={pageTitle} />
		</div>
		<Routes>
			<Route path='/' element={<HomePage />} />
		</Routes>
	</HashRouter>;
}

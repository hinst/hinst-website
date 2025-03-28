import { HashRouter, Route, Routes } from 'react-router';
import Header from './header';
import { useEffect, useState } from 'react';
import HomePage from './homePage';
import GoalBrowser from './personalGoals/goalBrowser';
import { getCurrentLanguage, SupportedLanguages } from './language';
import { DisplayWidthContext, LanguageContext } from './context';

const PAGE_TITLE = 'Showcase Website';

export default function App() {
	useEffect(() => {
		if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
			document.getElementsByTagName('html')[0].setAttribute('data-theme', 'dark');
	}, []);

	const [pageTitle, setPageTitle] = useState(PAGE_TITLE);
	useEffect(() => {
		document.title = pageTitle;
	}, [pageTitle]);

	const [currentLanguage, setCurrentLanguage] = useState<SupportedLanguages>(getCurrentLanguage());
	setInterval(() => {
		const newLanguage = getCurrentLanguage();
		if (newLanguage !== currentLanguage)
			setCurrentLanguage(newLanguage);
	}, 1000);

	const [displayWidth, setDisplayWidth] = useState(window.innerWidth);
	setInterval(() => {
		const newWidth = window.innerWidth;
		if (newWidth !== displayWidth)
			setDisplayWidth(newWidth);
	}, 1000);

	return <LanguageContext.Provider value={getCurrentLanguage()}>
		<DisplayWidthContext.Provider value={window.innerWidth}>
			<div
				style={{
					padding: 10,
					paddingBottom: 0,
					display: 'flex',
					flexDirection: 'column',
					width: '100%',
					maxWidth: '100%',
					maxHeight: '100%',
				}}
			>
				<HashRouter>
					<div style={{ marginBottom: 10 }}>
						<Header title={pageTitle} />
					</div>
					<Routes>
						<Route path='/' element={<HomePage setPageTitle={setPageTitle} />} />
						<Route path='/personal-goals/:id' element={<GoalBrowser setPageTitle={setPageTitle} />} />
					</Routes>
				</HashRouter>
			</div>
		</DisplayWidthContext.Provider>
	</LanguageContext.Provider>;
}

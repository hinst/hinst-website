import { CSSProperties, useEffect, useReducer } from 'react';
import { Info } from 'react-feather';
import { APP_TITLE, AUTHOR_NAME, COPYRIGHT_YEAR } from 'src/typescript/global';
import { SupportedLanguage, supportedLanguageNames } from 'src/typescript/language';
import { settingsStorage, Theme } from 'src/typescript/settings';

export default function SettingsPage(props: { setPageTitle: (title: string) => void }) {
	const [, forceUpdate] = useReducer((key) => key + 1, 0);

	useEffect(() => {
		props.setPageTitle('Settings');
	}, []);

	function setTheme(theme: Theme) {
		settingsStorage.theme = theme;
		forceUpdate();
	}

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
			<div>
				Color theme <br />
				<div className='ms-btn-group'>
					<ThemeButton
						theme={Theme.SYSTEM}
						currentTheme={settingsStorage.theme}
						title='System'
						setTheme={setTheme}
					/>
					<ThemeButton
						theme={Theme.DARK}
						currentTheme={settingsStorage.theme}
						title='Dark'
						setTheme={setTheme}
					/>
					<ThemeButton
						theme={Theme.LIGHT}
						currentTheme={settingsStorage.theme}
						title='Light'
						setTheme={setTheme}
					/>
				</div>
			</div>
			<div>
				Language <br />
				<div className='ms-btn-group'>
					<LanguageButton
						language={undefined}
						currentLanguage={settingsStorage.language}
						title='System'
						setLanguage={(language) => {
							settingsStorage.language = language;
							forceUpdate();
						}}
					/>
					{Object.values(SupportedLanguage).map((supportedLanguage) => (
						<LanguageButton
							key={supportedLanguage}
							language={supportedLanguage}
							currentLanguage={settingsStorage.language}
							title={supportedLanguageNames[supportedLanguage]}
							setLanguage={(language) => {
								settingsStorage.language = language;
								forceUpdate();
							}}
						/>
					))}
				</div>
				<br />
				<div
					style={{
						display: 'flex',
						alignItems: 'flex-start',
						gap: 5,
						marginTop: 4
					}}
				>
					<Info />
					Translations are currently available only for articles in personal goals. <br />
					Buttons and other UI elements are English only.
				</div>
			</div>
			<div>
				{APP_TITLE} &copy; {COPYRIGHT_YEAR} {AUTHOR_NAME}
			</div>
		</div>
	);
}

const activeStyle: CSSProperties = {
	backgroundColor: 'rgba(var(--primary-bg-color), 1)',
	color: 'white'
};

function ThemeButton(props: {
	theme: Theme;
	currentTheme: Theme;
	title: string;
	setTheme: (theme: Theme) => void;
}) {
	const style = props.currentTheme === props.theme ? activeStyle : {};
	return (
		<button
			type='button'
			className='ms-btn ms-outline'
			style={style}
			onClick={() => props.setTheme(props.theme)}
		>
			{props.title}
		</button>
	);
}

function LanguageButton(props: {
	language: SupportedLanguage | undefined;
	currentLanguage: SupportedLanguage | undefined;
	title: string;
	setLanguage: (language: SupportedLanguage | undefined) => void;
}) {
	const style = props.currentLanguage === props.language ? activeStyle : {};
	return (
		<button
			type='button'
			className='ms-btn ms-outline'
			style={style}
			onClick={() => props.setLanguage(props.language)}
		>
			{props.title}
		</button>
	);
}

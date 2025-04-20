import { CSSProperties, useEffect, useReducer } from 'react';
import { APP_TITLE, AUTHOR_NAME, COPYRIGHT_YEAR } from 'src/typescript/global';
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
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<div>
				Color theme &nbsp;
				<div className='ms-btn-group'>
					<CreateThemeButton
						theme={Theme.SYSTEM}
						currentTheme={settingsStorage.theme}
						title='System'
						setTheme={setTheme}
					/>
					<CreateThemeButton
						theme={Theme.DARK}
						currentTheme={settingsStorage.theme}
						title='Dark'
						setTheme={setTheme}
					/>
					<CreateThemeButton
						theme={Theme.LIGHT}
						currentTheme={settingsStorage.theme}
						title='Light'
						setTheme={setTheme}
					/>
				</div>
			</div>
			<div>
				{APP_TITLE} &copy; {COPYRIGHT_YEAR} {AUTHOR_NAME}
			</div>
		</div>
	);
}

function CreateThemeButton(props: {
	theme: Theme;
	currentTheme: Theme;
	title: string;
	setTheme: (theme: Theme) => void;
}) {
	const style: CSSProperties = {};
	if (props.currentTheme === props.theme) {
		style.backgroundColor = 'rgba(var(--primary-bg-color), 1)';
		style.color = 'white';
	}
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

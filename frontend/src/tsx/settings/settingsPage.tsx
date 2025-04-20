import { useEffect, useReducer } from 'react';
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

	function getButtonClass(theme: Theme) {
		return settingsStorage.theme === theme ? 'ms-primary' : 'ms-outline';
	}

	return (
		<div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
			<div>
				Color theme &nbsp;
				<div className='ms-btn-group'>
					<button
						type='button'
						className={'ms-btn ' + getButtonClass(Theme.SYSTEM)}
						onClick={() => setTheme(Theme.SYSTEM)}
					>
						System
					</button>
					<button
						type='button'
						className={'ms-btn ' + getButtonClass(Theme.DARK)}
						onClick={() => setTheme(Theme.DARK)}
					>
						Dark
					</button>
					<button
						type='button'
						className={'ms-btn ' + getButtonClass(Theme.LIGHT)}
						onClick={() => setTheme(Theme.LIGHT)}
					>
						Light
					</button>
				</div>
			</div>
			<div>
				{APP_TITLE} &copy; {COPYRIGHT_YEAR} {AUTHOR_NAME}
			</div>
		</div>
	);
}

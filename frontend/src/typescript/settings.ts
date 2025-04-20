export enum Theme {
	LIGHT = 'LIGHT',
	DARK = 'DARK',
	SYSTEM = 'SYSTEM',
}

class SettingsStorage {
	initialize() {
		this.applyTheme();
	}

	private applyTheme() {
		if (this.resolvedTheme === Theme.DARK)
			document.documentElement.setAttribute('data-theme', 'dark');
		else
			document.documentElement.removeAttribute('data-theme');
	}

	get theme(): Theme {
		const theme = localStorage.getItem('theme');
		if (theme) return theme as Theme;
		return Theme.SYSTEM;
	}

	get resolvedTheme(): Theme {
		const theme = this.theme;
		if (theme === Theme.SYSTEM) {
			if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
				return Theme.DARK;
			} else {
				return Theme.LIGHT;
			}
		}
		return theme;
	}

	set theme(value: Theme) {
		if (value === Theme.SYSTEM)
			localStorage.removeItem('theme');
		else
			localStorage.setItem('theme', value);
		this.applyTheme();
	}
}

export const settingsStorage = new SettingsStorage();
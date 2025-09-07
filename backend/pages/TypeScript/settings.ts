export enum Theme {
	LIGHT = 'LIGHT',
	DARK = 'DARK',
	SYSTEM = 'SYSTEM',
}

class SettingsStorage {
	initialize() {
		this.applyTheme();
	}

	private readonly keyTheme = 'theme';

	get theme(): Theme {
		const theme = localStorage.getItem(this.keyTheme);
		if (theme) return theme as Theme;
		return Theme.SYSTEM;
	}

	set theme(value: Theme) {
		if (value === Theme.SYSTEM)
			localStorage.removeItem(this.keyTheme);
		else
			localStorage.setItem(this.keyTheme, value);
		this.applyTheme();
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

	private applyTheme() {
		const dataThemeAttributeKey = 'data-theme';
		if (this.resolvedTheme === Theme.DARK)
			document.documentElement.setAttribute(dataThemeAttributeKey, 'dark');
		else
			document.documentElement.removeAttribute(dataThemeAttributeKey);
	}
}

export const settingsStorage = new SettingsStorage();
const headerLanguageButton = document.getElementById('hamburger-menu-button') as HTMLButtonElement;
const headerLanguagePopup = document.getElementById('hamburger-menu-popup') as HTMLDivElement;

export function main() {
	if (headerLanguageButton && headerLanguagePopup) {
		const displayStyle = 'flex';
		headerLanguageButton.onclick = () => {
			const isVisible = headerLanguagePopup.style.display === displayStyle;
			headerLanguagePopup.style.display = isVisible ? 'none' : displayStyle;
		};
		const links = headerLanguagePopup.getElementsByTagName('a');
		let lastLink: HTMLAnchorElement | null = null;
		for (const link of links) {
			const href = link.getAttribute('href') || '';
			if (window.location.pathname.startsWith(href))
				lastLink = link;
		}
		if (lastLink)
			lastLink.classList.remove('ms-outline');
	}
}
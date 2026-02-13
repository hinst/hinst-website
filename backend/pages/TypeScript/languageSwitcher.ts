const hamburgerMenuButton = document.getElementById('hamburger-menu-button') as HTMLButtonElement;
const hamburgerMenuPopup = document.getElementById('hamburger-menu-popup') as HTMLDivElement;

export function main() {
	if (hamburgerMenuButton && hamburgerMenuPopup) {
		const displayStyle = 'flex';
		hamburgerMenuButton.onclick = () => {
			const isVisible = hamburgerMenuPopup.style.display === displayStyle;
			hamburgerMenuPopup.style.display = isVisible ? 'none' : displayStyle;
		};
		const links = hamburgerMenuPopup.getElementsByTagName('a');
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
import * as goalPost from './goalPost';
goalPost;

const headerLanguageButton = document.getElementById('header-language-button') as HTMLButtonElement;
const headerLanguagePopup = document.getElementById('header-language-popup') as HTMLDivElement;
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
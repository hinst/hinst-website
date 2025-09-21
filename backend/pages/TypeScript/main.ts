import * as goalPost from './goalPost';
goalPost;

const headerLanguageButton = document.getElementById('header-language-button') as HTMLButtonElement;
const headerLanguagePopup = document.getElementById('header-language-popup') as HTMLDivElement;
if (headerLanguageButton && headerLanguagePopup) {
	const displayStyle = 'flex';
	headerLanguageButton.onclick = () => {
		headerLanguagePopup.style.display = headerLanguagePopup.style.display === displayStyle ? 'none' : displayStyle;
	};
}
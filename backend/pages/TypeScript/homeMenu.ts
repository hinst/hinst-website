const hamburgerMenuButton = document.getElementById('hamburger-menu-button') as HTMLButtonElement;
const hamburgerMenuButtonInner = document.getElementById(
	'hamburger-menu-button-inner'
) as HTMLDivElement;
const hamburgerMenuPopup = document.getElementById('hamburger-menu-popup') as HTMLDivElement;

let hamburgerMenuVisible = false;

function updateHamburgerMenuVisibility() {
	if (null == hamburgerMenuButton || null == hamburgerMenuPopup) {
		console.error('Hamburger menu elements are missing found');
		return;
	}
	if (hamburgerMenuVisible) {
		hamburgerMenuPopup.style.transform = 'translate(0,0)';
		hamburgerMenuPopup.style.opacity = '1';
		hamburgerMenuPopup.classList.add('ms-box-shadow');
	} else {
		hamburgerMenuPopup.style.transform = 'translate(-100%, 0)';
		hamburgerMenuPopup.style.opacity = '0';
		hamburgerMenuPopup.classList.remove('ms-box-shadow');
	}
}

function initializeHamburgerMenu() {
	function toggle() {
		console.log('Hamburger menu button clicked');
		hamburgerMenuVisible = !hamburgerMenuVisible;
		updateHamburgerMenuVisibility();
	}
	hamburgerMenuButton.onclick = toggle;
	hamburgerMenuButtonInner.onclick = toggle;
}

function highlightCurrentLanguage() {
	const links = hamburgerMenuPopup.getElementsByTagName('a');
	let lastLink: HTMLAnchorElement | null = null;
	for (const link of links) {
		const href = link.getAttribute('href') || '';
		if (window.location.pathname.startsWith(href)) lastLink = link;
	}
	if (lastLink) {
		lastLink.classList.remove('ms-outline');
		lastLink.classList.add('ms-primary');
	}
}

export function main() {
	if (hamburgerMenuButton && hamburgerMenuPopup) {
		initializeHamburgerMenu();
		highlightCurrentLanguage();
	}
}

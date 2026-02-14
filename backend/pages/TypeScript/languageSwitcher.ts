const hamburgerMenuButton = document.getElementById('hamburger-menu-button') as HTMLButtonElement;
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

export function main() {
	if (hamburgerMenuButton && hamburgerMenuPopup) {
		hamburgerMenuButton.onclick = () => {
			console.log('Hamburger menu button clicked');
			hamburgerMenuVisible = !hamburgerMenuVisible;
			updateHamburgerMenuVisibility();
		};
		const links = hamburgerMenuPopup.getElementsByTagName('a');
		let lastLink: HTMLAnchorElement | null = null;
		for (const link of links) {
			const href = link.getAttribute('href') || '';
			if (window.location.pathname.startsWith(href)) lastLink = link;
		}
		if (lastLink) lastLink.classList.remove('ms-outline');
	}
}

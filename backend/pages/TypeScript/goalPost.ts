let background: string | undefined = undefined;
function updateSize() {
	const outerElement = document.querySelector('.goal-post-outer') as HTMLElement;
	const innerElement = document.querySelector('.goal-post-inner') as HTMLElement;
	if (outerElement && innerElement) {
		if (!background)
			background = outerElement.style.background;
		const isBackgroundNeeded = outerElement.clientWidth - innerElement.clientWidth >= 100;
		outerElement.style.background = isBackgroundNeeded ? background : 'none';
	}
}
updateSize();
setInterval(updateSize, 500);
export const goalPost = 1;
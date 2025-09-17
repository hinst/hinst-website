function main() {
	const outerElement = document.querySelector('.goal-post-outer') as HTMLElement;
	if (!outerElement)
		return;
	const innerElement = document.querySelector('.goal-post-inner') as HTMLElement;
	if (!innerElement)
		return;
	const background = outerElement.style.background;

	function updateSize() {
		const isBackgroundNeeded = outerElement.clientWidth - innerElement.clientWidth >= 100;
		outerElement.style.background = isBackgroundNeeded ? background : 'none';
	}
	updateSize();
	new ResizeObserver(updateSize).observe(outerElement);
}
main();

export const goalPost = 1;
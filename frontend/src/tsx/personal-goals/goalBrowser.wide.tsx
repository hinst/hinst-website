import type { ReactElement } from 'react';
import { useEffect, useRef, useState } from 'react';

const ARTICLE_PADDING = 20;
const ARTICLE_WIDTH = 1000 + ARTICLE_PADDING * 2;
const STRIPES_BACKGROUND = `repeating-linear-gradient(
	45deg,
	rgba(var(--main-bg), 1),
	rgba(var(--main-bg), 1) 10px,
	rgba(var(--light-bg-color), 1) 10px,
	rgba(var(--light-bg-color), 1) 20px
)`;
const STRIPES_MIN_WIDTH = 100;

export default function GoalBrowserWide(props: {
	activePostDate: string;
	getGoalCalendarPanel: () => ReactElement;
	getGoalPostPanel: () => ReactElement;
}) {
	const articleContainerRef = useRef<HTMLDivElement>(null);
	const [articleContainerWidth, setArticleContainerWidth] = useState(0);

	useEffect(() => {
		const element = articleContainerRef.current;
		if (!element) return;
		const ro = new ResizeObserver(([entry]) =>
			setArticleContainerWidth(entry.contentRect.width)
		);
		ro.observe(element);
		return () => ro.disconnect();
	}, []);

	return (
		<div
			style={{
				display: 'flex',
				gap: 20,
				minHeight: 0,
				height: '100%'
			}}
		>
			<div
				style={{
					display: 'flex',
					overflowY: 'auto',
					flexShrink: 0,
					flexBasis: 'fit-content'
				}}
			>
				{props.getGoalCalendarPanel()}
			</div>
			<div
				ref={articleContainerRef}
				style={{
					flexGrow: 1,
					justifyContent: 'center',
					display: 'flex',
					minHeight: 0,
					maxHeight: '100%',
					background:
						articleContainerWidth > ARTICLE_WIDTH + STRIPES_MIN_WIDTH
							? STRIPES_BACKGROUND
							: undefined
				}}
			>
				<div
					className='ms-bg-main'
					style={{
						paddingLeft: ARTICLE_PADDING,
						paddingRight: ARTICLE_PADDING,
						flexGrow: 1,
						maxWidth: ARTICLE_WIDTH,
						backgroundAttachment: 'fixed',
						minHeight: 0,
						overflowY: 'auto'
					}}
				>
					{props.activePostDate ? props.getGoalPostPanel() : undefined}
				</div>
			</div>
		</div>
	);
}

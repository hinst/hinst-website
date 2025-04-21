import { GoalPostObjectExtended } from 'src/typescript/personal-goals/smartPost';
import SafeHtmlView from '../safeHtmlView';
import { Info } from 'react-feather';
import { getHashFromString } from 'src/typescript/string';

export default function GoalPostView(props: { postData: GoalPostObjectExtended }) {
	return (
		<div>
			{props.postData.isAutoTranslated ? (
				<div
					className='ms-alert ms-light'
					style={{ display: 'flex', alignItems: 'center', gap: 8 }}
				>
					<div>
						<Info />
					</div>
					<div>
						This text was automatically translated to {props.postData.languageName}{' '}
						language by LLM tool.
					</div>
				</div>
			) : undefined}
			{props.postData.isTranslationPending ? (
				<div
					className='ms-alert ms-light'
					style={{ display: 'flex', alignItems: 'center', gap: 8 }}
				>
					<div>
						<Info className='ms-text-secondary' />
					</div>
					<div>
						The automatic translation of this text to {props.postData.languageName}
						language is not available yet. Please come back later or check older posts.
					</div>
				</div>
			) : undefined}
			<div className='goalPostViewText'>
				<SafeHtmlView htmlText={props.postData.text} updateDocument={removeRedirect} />
			</div>
			<div
				style={{
					display: 'flex',
					flexWrap: 'wrap',
					gap: 10
				}}
			>
				{props.postData.images?.map((image) => (
					<GoalImage key={getHashFromString(image)} data={image} />
				))}
			</div>
		</div>
	);
}

function GoalImage(props: { data: string }) {
	return (
		<a href={props.data}>
			<img
				className='ms-card ms-border'
				width={240}
				height={240}
				style={{
					width: 240,
					height: 240,
					objectFit: 'cover',
					margin: 0,
					padding: 0
				}}
				src={props.data}
				alt='Image'
			/>
		</a>
	);
}

function removeRedirect(document: Document) {
	const prefix = 'http://smartprogress.do/site/redirect/?url=';
	[...document.getElementsByTagName('a')].forEach((link) => {
		try {
			let href = link.getAttribute('href');
			if (href?.startsWith(prefix)) {
				href = href.substring(prefix.length);
				if (href.endsWith('%'))
					href = href.substring(0, href.length - 1);
				href = decodeURIComponent(href);
				link.setAttribute('href', href);
			}
		} catch (e) {
			console.warn('Cannot process link', link, e);
		}
	});
}

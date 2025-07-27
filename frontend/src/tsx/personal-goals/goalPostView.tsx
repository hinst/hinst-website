import { GoalPostObject } from 'src/typescript/personal-goals/smartPost';
import SafeHtmlView from '../safeHtmlView';
import { apiClient } from 'src/typescript/apiClient';
import { Info } from 'react-feather';
import { getHashFromString } from 'src/typescript/string';

export default function GoalPostView(props: { postData: GoalPostObject }) {
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
				{new Array(props.postData.imageCount).fill(undefined).map((_, index) => (
					<GoalImage
						key={getHashFromString('' + props.postData.goalId + index)}
						goalId={props.postData.goalId}
						postDateTime={props.postData.dateTime}
						index={index}
					/>
				))}
			</div>
		</div>
	);
}

function GoalImage(props: { goalId: number; postDateTime: number; index: number }) {
	const url = apiClient.getImageUrl(props.goalId, props.postDateTime, props.index);
	return (
		<a href={url}>
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
				src={url}
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
				if (href.endsWith('%')) href = href.substring(0, href.length - 1);
				href = decodeURIComponent(href);
				link.setAttribute('href', href);
			}
		} catch (e) {
			console.warn('Cannot process link', link, e);
		}
	});
}

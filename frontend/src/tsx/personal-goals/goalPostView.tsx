import { GoalPostObjectExtended } from 'src/typescript/personal-goals/smartPost';
import SafeHtmlView from '../safeHtmlView';
import { Info } from 'react-feather';

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
			{props.postData.languageNamePending ? (
				<div
					className='ms-alert ms-light'
					style={{ display: 'flex', alignItems: 'center', gap: 8 }}
				>
					<div>
						<Info className='ms-text-secondary' />
					</div>
					<div>
						The automatic translation of this text to{' '}
						{props.postData.languageNamePending} language is not available yet. Please
						come back later or check older posts.
					</div>
				</div>
			) : undefined}
			<div className='goalPostViewText'>
				<SafeHtmlView htmlText={props.postData.text} />
			</div>
			<div
				style={{
					display: 'flex',
					flexWrap: 'wrap',
					gap: 10
				}}
			>
				{props.postData.images?.map((image) => <GoalImage key={image.slice(0, 100)} data={image} />)}
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

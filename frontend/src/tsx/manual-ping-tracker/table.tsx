import { DateTime } from 'luxon';
import { CSSProperties, useState } from 'react';
import { Check, CheckCircle, Copy } from 'react-feather';
import { apiClient } from 'src/typescript/apiClient';
import type { GoalPostSearchIndexingHeader } from 'src/typescript/apiTypes';

const BUTTON_STYLE: CSSProperties = {
	display: 'flex',
	alignItems: 'center',
	justifyContent: 'center',
	padding: '6px 12px',
	height: 38,
	width: 160
};

export function HeaderRow() {
	return (
		<tr>
			<th>URL</th>
			<th>Google Ping</th>
		</tr>
	);
}

export function Row(props: { record: GoalPostSearchIndexingHeader; onPinged: () => void }) {
	const [isCopied, setIsCopied] = useState(false);
	const [isPinged, setIsPinged] = useState(false);
	return (
		<tr>
			<td>
				<div
					style={{
						display: 'flex',
						flexDirection: 'column',
						gap: 8,
						alignItems: 'flex-start',
						width: 'fit-content'
					}}
				>
					<span>{props.record.title}</span>
					<a href={props.record.publicUrl}>{props.record.publicUrl}</a>
				</div>
			</td>
			<td>
				<div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
					<div>
						Pinged at{' '}
						<code>
							{' '}
							{props.record.googlePingedAt
								? formatDate(props.record.googlePingedAt)
								: 'never'}
						</code>
					</div>
					<div>
						Indexing status:{' '}
						<code>{props.record.googleSearchIndexingStatus || '?'}</code>
					</div>
					{isPinged ? (
						<button
							type='button'
							className='ms-btn ms-outline'
							style={BUTTON_STYLE}
							disabled={true}
						>
							<CheckCircle /> &nbsp; Done
						</button>
					) : isCopied ? (
						<PingUrlButton
							onDone={() => {
								setIsPinged(true);
								props.onPinged();
							}}
							goalId={props.record.goalId}
							postDateTime={props.record.dateTime}
						/>
					) : (
						<CopyUrlButton
							onDone={() => setIsCopied(true)}
							url={'' + props.record.publicUrl}
						/>
					)}
				</div>
			</td>
		</tr>
	);
}

function formatDate(timestamp: number | null) {
	if (timestamp == null) {
		return '';
	}
	const date = DateTime.fromMillis(timestamp * 1000);
	return date.toFormat('yyyy-MM-dd');
}

function CopyUrlButton(props: { onDone: () => void; url: string }) {
	function copyUrl() {
		let ok = false;
		try {
			navigator.clipboard.writeText(props.url);
			ok = true;
		} catch (error) {
			const message = 'Cannot copy URL to clipboard';
			alert(message);
			console.error(message, error);
		}
		if (ok) props.onDone();
	}
	return (
		<button
			title='Copy URL to the clipboard. You should paste it into Google Search Console and click Request Indexing'
			type='button'
			className='ms-btn ms-action'
			onClick={copyUrl}
			style={BUTTON_STYLE}
		>
			<Copy />
			&nbsp; Copy URL
		</button>
	);
}

function PingUrlButton(props: { onDone: () => void; goalId: number; postDateTime: number }) {
	const [isLoading, setIsLoading] = useState(false);
	async function ping() {
		setIsLoading(true);
		await apiClient.setGoalPostGooglePingedAt(props.goalId, props.postDateTime);
		setIsLoading(false);
		props.onDone();
	}
	return (
		<button
			title='Confirm that the URL was accepted by Google Search Console'
			type='button'
			className='ms-btn ms-action2'
			onClick={ping}
			disabled={isLoading}
			style={BUTTON_STYLE}
		>
			<Check />
			&nbsp; Mark pinged
		</button>
	);
}

import { DateTime } from 'luxon';
import { useState } from 'react';
import { Check, CheckCircle, Copy } from 'react-feather';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';

export function Row(props: { record: UrlPingRecord }) {
	const [isCopied, setIsCopied] = useState(false);
	const [isPinged, setIsPinged] = useState(false);
	return (
		<tr>
			<td>{props.record.url}</td>
			<td>{formatDate(props.record.googlePingedAt)}</td>
			<td style={{ height: 62 }}>
				{isPinged ? (
					<div style={{ display: 'flex', alignItems: 'center' }}>
						<CheckCircle /> &nbsp; Done
					</div>
				) : isCopied ? (
					<PingUrlButton onDone={() => setIsPinged(true)} url={props.record.url} />
				) : props.record.googlePingedManuallyAt != null ? (
					formatDate(props.record.googlePingedManuallyAt)
				) : (
					<CopyUrlButton onDone={() => setIsCopied(true)} url={props.record.url} />
				)}
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
			style={{
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'center',
				padding: '6px 12px',
				width: 160
			}}
		>
			<Copy />
			&nbsp; Copy URL
		</button>
	);
}

function PingUrlButton(props: { onDone: () => void; url: string }) {
	const [isLoading, setIsLoading] = useState(false);
	async function ping() {
		setIsLoading(true);
		await apiClient.pingUrlManually(props.url);
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
			style={{
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'center',
				padding: '6px 12px',
				width: 160
			}}
		>
			<Check />
			&nbsp; Confirm pinged
		</button>
	);
}

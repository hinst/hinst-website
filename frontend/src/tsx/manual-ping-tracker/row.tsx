import { DateTime } from 'luxon';
import { useState } from 'react';
import { Check, Copy } from 'react-feather';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';

export function Row(props: { record: UrlPingRecord }) {
	const [isLoading, setIsLoading] = useState(false);
	const [isDone, setIsDone] = useState(false);
	async function pingNow() {
		setIsLoading(true);
		{
			navigator.clipboard.writeText(props.record.url);
			await apiClient.pingUrlManually(props.record.url);
			setIsDone(true);
		}
		setIsLoading(false);
	}
	return (
		<tr>
			<td>{props.record.url}</td>
			<td>{formatDate(props.record.googlePingedAt)}</td>
			<td style={{ height: 62 }}>
				{isDone ? (
					<div style={{ display: 'flex', alignItems: 'center' }}>
						<Check /> &nbsp; Done
					</div>
				) : props.record.googlePingedManuallyAt != null ? (
					formatDate(props.record.googlePingedManuallyAt)
				) : (
					<button
						title='Copy URL to the clipboard and mark the URL as pinged. You gotta paste it into Google Search Console yourself'
						type='button'
						className='ms-btn ms-action'
						onClick={pingNow}
						disabled={isLoading}
						style={{ display: 'flex', alignItems: 'center', padding: '6px 12px' }}
					>
						<Copy />
						&nbsp; Commit
					</button>
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

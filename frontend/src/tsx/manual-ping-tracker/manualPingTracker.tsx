import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';
import { DateTime } from 'luxon';
import { Check, Copy } from 'react-feather';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	const [urlPings, setUrlPings] = useState<Array<UrlPingRecord>>([]);
	const [manuallyPingedVisible, setManuallyPingedVisible] = useState(true);
	async function loadUrlPings() {
		const urlPings = await apiClient.getUrlPings();
		validatePings(urlPings);
		setUrlPings(urlPings);
	}
	function getVisibleUrlPings() {
		if (manuallyPingedVisible) {
			return urlPings;
		}
		return urlPings.filter((ping) => ping.googlePingedManuallyAt == null);
	}
	useEffect(() => {
		props.setPageTitle('Manual URL ping tracker');
		loadUrlPings();
	}, []);
	return (
		<div>
			<div
				className='ms-alert ms-light'
				style={{ display: 'flex', alignItems: 'center', gap: 10, padding: 8 }}
			>
				<div>Showing URLs: {getVisibleUrlPings().length}</div>
				<b>|</b>
				<div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
					<input
						type='checkbox'
						checked={manuallyPingedVisible}
						onChange={() => setManuallyPingedVisible(!manuallyPingedVisible)}
					/>
					Manually pinged
				</div>
			</div>
			<table className='ms-table ms-striped'>
				<thead>
					<tr>
						<th>URL</th>
						<th>Google Ping</th>
						<th>Manual Ping</th>
					</tr>
				</thead>
				<tbody>
					{getVisibleUrlPings().map((item) => (
						<Row key={item.url} record={item} />
					))}
				</tbody>
			</table>
		</div>
	);
}

function Row(props: { record: UrlPingRecord }) {
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

function validatePings(urlPings: UrlPingRecord[]) {
	const urls = new Set<string>();
	for (const ping of urlPings)
		if (urls.has(ping.url)) console.warn('Duplicate URL ping record for URL: ' + ping.url);
}

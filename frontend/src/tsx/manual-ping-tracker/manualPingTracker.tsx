import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';
import { DateTime } from 'luxon';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	const [urlPings, setUrlPings] = useState<Array<UrlPingRecord>>([]);
	async function loadUrlPings() {
		const urlPings = await apiClient.getUrlPings();
		validatePings(urlPings);
		setUrlPings(urlPings);
	}
	useEffect(() => {
		props.setPageTitle('Manual URL ping tracker');
		loadUrlPings();
	}, []);
	return (
		<table className='ms-table ms-striped'>
			<thead>
				<tr>
					<th>URL</th>
					<th>Google Ping</th>
					<th>Manual Ping</th>
				</tr>
			</thead>
			<tbody>
				{urlPings.map((item) => (
					<Row key={item.url} record={item} />
				))}
			</tbody>
		</table>
	);
}

function Row(props: { record: UrlPingRecord }) {
	const [isDone, setIsDone] = useState(false);
	async function pingNow() {
		navigator.clipboard.writeText(props.record.url);
		await apiClient.pingUrlManually(props.record.url);
	}
	return (
		<tr>
			<td>{props.record.url}</td>
			<td>{formatDate(props.record.googlePingedAt)}</td>
			<td>
				{isDone ? 'Done' : ''}
				{props.record.googlePingedManuallyAt != null ? (
					formatDate(props.record.googlePingedManuallyAt)
				) : (
					<button
						title='Copy URL to the clipboard and mark the URL as pinged. You gotta paste it into Google Search Console yourself'
						type='button'
						className='ms-btn ms-action'
						onClick={pingNow}
					>
						Commit
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
	for (const ping of urlPings) {
		if (urls.has(ping.url)) console.warn('Duplicate URL ping record for URL: ' + ping.url);
	}
}

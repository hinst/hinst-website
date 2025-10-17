import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';
import { Row } from './row';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	const [urlPings, setUrlPings] = useState<Array<UrlPingRecord>>([]);
	const [manuallyPingedVisible, setManuallyPingedVisible] = useState(true);
	async function loadUrlPings() {
		const urlPings = await apiClient.getUrlPings();
		validatePings(urlPings);
		setUrlPings(urlPings);
	}
	function getVisibleUrlPings() {
		if (manuallyPingedVisible) return urlPings;
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

function validatePings(urlPings: UrlPingRecord[]) {
	const urls = new Set<string>();
	for (const ping of urlPings)
		if (urls.has(ping.url)) console.warn('Duplicate URL ping record for URL: ' + ping.url);
}

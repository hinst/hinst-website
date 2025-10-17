import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { UrlPingRecord } from 'src/typescript/urlPing';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	const [urlPings, setUrlPings] = useState<Array<UrlPingRecord>>([]);
	async function loadUrlPings() {
		const urlPings = await apiClient.getUrlPings();
		setUrlPings(urlPings);
	}
	useEffect(() => {
		props.setPageTitle('Manual URL ping tracker');
		loadUrlPings();
	}, []);
	return <div>{urlPings.map((item) => Row(item))}</div>;
}

function Row(record: UrlPingRecord) {
	return <div>{record.url}</div>;
}

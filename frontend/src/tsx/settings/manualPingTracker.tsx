import { useEffect } from 'react';
import { apiClient } from 'src/typescript/apiClient';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	async function loadUrlPings() {
		const urlPings = apiClient.getUrlPings();
		console.log({ urlPings });
	}
	useEffect(() => {
		props.setPageTitle('Manual URL ping tracker');
		loadUrlPings();
	}, []);
	return <div></div>;
}

import { useEffect } from 'react';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	useEffect(() => {
		props.setPageTitle('Manual URL ping tracker');
	}, []);
}

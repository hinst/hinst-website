import { useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import { HeaderRow, Row } from './table';
import type { GoalPostSearchIndexingHeader } from 'src/typescript/apiTypes';

export default function ManualPingTracker(props: { setPageTitle: (title: string) => void }) {
	const [urlPings, setUrlPings] = useState<Array<GoalPostSearchIndexingHeader>>([]);
	const [manuallyPingedVisible, setManuallyPingedVisible] = useState(true);
	async function loadUrlPings() {
		const urlPings = await apiClient.getUrlPings();
		setUrlPings(urlPings);
	}
	function getVisibleUrlPings() {
		return urlPings;
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
					<HeaderRow />
				</thead>
				<tbody>
					{getVisibleUrlPings().map((item) => (
						<Row key={item.publicUrl} record={item} onPinged={loadUrlPings} />
					))}
				</tbody>
			</table>
		</div>
	);
}

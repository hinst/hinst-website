import { useContext, useEffect, useState } from 'react';
import { apiClient } from 'src/typescript/apiClient';
import type { GoalPostSearchIndexingHeader } from 'src/typescript/apiTypes';
import { AppContext } from 'src/typescript/appContext';
import { PageTitle } from 'src/typescript/pageTitle';
import { HeaderRow, Row } from './table';

export default function ManualPingTracker() {
	const context = useContext(AppContext);
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
		context.setPageTitle(new PageTitle('Administrator', 'Manual URL ping tracker'));
		const _promise = loadUrlPings();
	}, []);
	return (
		<div>
			<div
				className='ms-alert ms-light'
				style={{ display: 'none', alignItems: 'center', gap: 10, padding: 8 }}
			>
				<div>
					Showing URLs: {getVisibleUrlPings().length} of {urlPings.length}
				</div>
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

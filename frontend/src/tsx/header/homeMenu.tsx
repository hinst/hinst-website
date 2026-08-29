import { Activity, Search, Settings } from 'react-feather';
import { NavLink } from 'react-router';
import { useEffect, useState } from 'react';
import { apiClient } from '../../typescript/apiClient';

export function HomeMenu() {
	const [isAdminModeEnabled, setIsAdminModeEnabled] = useState(false);

	useEffect(() => {
		apiClient.isAdminModeEnabled().then(setIsAdminModeEnabled);
	}, []);

	return (
		<div style={{ display: 'flex', gap: 5, flexDirection: 'column', width: 160 }}>
			<NavLink
				to='/settings'
				className='ms-btn ms-outline'
				style={{
					display: 'flex',
					alignItems: 'center',
					gap: 5,
					padding: 5,
					paddingRight: 10
				}}
			>
				<Settings /> Settings
			</NavLink>
			<NavLink
				to='/personal-goals-search'
				className='ms-btn ms-outline'
				style={{
					display: 'flex',
					alignItems: 'center',
					gap: 5,
					padding: 5,
					paddingRight: 10
				}}
			>
				<Search /> Search
			</NavLink>
			{isAdminModeEnabled && (
				<NavLink
					to='/manual-ping-tracker'
					className='ms-btn ms-outline'
					style={{
						display: 'flex',
						alignItems: 'center',
						gap: 5,
						padding: 5,
						paddingRight: 10
					}}
				>
					<Activity /> Ping Tracker
				</NavLink>
			)}
		</div>
	);
}

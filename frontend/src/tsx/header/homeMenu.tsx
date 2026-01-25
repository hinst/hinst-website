import { Search, Settings } from 'react-feather';
import { NavLink } from 'react-router';

export function HomeMenu() {
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
		</div>
	);
}

import { Settings } from 'react-feather';
import { NavLink } from 'react-router';

export function HomeMenu() {
	return (
		<div>
			<NavLink
				to='/settings'
				type='button'
				className='ms-btn ms-outline'
				style={{ display: 'flex', alignItems: 'center', gap: 5, padding: 5 }}
			>
				<Settings /> Settings
			</NavLink>
		</div>
	);
}

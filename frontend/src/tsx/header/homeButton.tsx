// @ts-ignore
import icon from 'images/icon.webp';
import { NavLink } from 'react-router';

export function HomeButton() {
	return (
		<NavLink to='/'>
			<img
				src={icon}
				width={42}
				height={42}
				style={{
					width: 42,
					height: 42,
					borderRadius: '50%'
				}}
				alt='icon'
				className='hover-outline'
			/>
		</NavLink>
	);
}

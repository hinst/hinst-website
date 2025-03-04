import { NavLink, useLocation } from 'react-router';
// @ts-ignore
import icon from './icon.webp';

export default function Header(props: { title: string }) {
	const location = useLocation();
	return <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
		<NavLink to='/'>
		<img
			src={icon}
			width={42}
			height={42}
			style={{ borderRadius: '50%' }}
			alt='icon'
			className=''
		/>
		</NavLink>
		<h6 style={{marginTop: 10, marginBottom: 10}}>{props.title}</h6>
	</div>;
}
import { NavLink, useLocation } from 'react-router';
// @ts-ignore
import icon from './icon.webp';
import { useState } from 'react';
import { Info } from 'react-feather';
import { AUTHOR_NAME, COPYRIGHT_YEAR } from './global';

export default function Header(props: { title: string }) {
	const [isCopyrightVisible, setCopyrightVisible] = useState(false);
	const location = useLocation();
	return <div style={{ display: 'flex', alignItems: 'center', gap: 10, maxWidth: '100%' }}>
		<NavLink to='/'>
			<img
				src={icon}
				width={42}
				height={42}
				style={{ borderRadius: '50%' }}
				alt='icon'
				className='hover-outline'
			/>
		</NavLink>
		<h6
			style={{
				marginTop: 10,
				marginBottom: 10,
				textWrap: 'nowrap',
				overflowY: 'clip',
				textOverflow: 'ellipsis',
				flexShrink: 0,
				flexBasis: 0,
				minWidth: 0,
			}}
		>
			{props.title}
		</h6>

		<div style={{flexGrow: 1}}></div>

		<div
			style={{
				display: 'flex',
				alignItems: 'center',
				flexShrink: 1,
				gap: 10,
			}}
			className='ms-bg-main blurry-main-background'
		>
			<div
				className='ms-bg-main'
				style={{
					display: isCopyrightVisible ? 'block' : 'none',
					textWrap: 'nowrap',
					marginLeft: 5,
				}}
			>
				&copy; {COPYRIGHT_YEAR} {AUTHOR_NAME}
			</div>
			<button
				onClick={() => setCopyrightVisible(!isCopyrightVisible)}
				className='ms-btn ms-rounded'
				style={{width: 40, height: 40, position: 'relative'}}
			>
				<Info
					size={20}
					style={{
						position: 'absolute',
						left: '50%',
						top: '50%',
						transform: 'translate(-50%, -50%)',
					}}
				/>
			</button>
		</div>
	</div>;
}
import { NavLink } from 'react-router';
// @ts-ignore
import icon from 'images/icon.webp';
import { Settings } from 'react-feather';
import { AUTHOR_NAME, COPYRIGHT_YEAR } from 'src/typescript/global';

export default function Header(props: { title: string }) {
	return (
		<div style={{ display: 'flex', alignItems: 'center', gap: 10, maxWidth: '100%' }}>
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
			<div
				style={{
					display: 'flex',
					flexDirection: 'column',
					overflowY: 'clip',
					gap: 4,
					flexShrink: 0,
					flexBasis: 0,
					minWidth: 0
				}}
			>
				<b
					style={{
						textWrap: 'nowrap',
						textOverflow: 'ellipsis',
						overflowY: 'clip'
					}}
				>
					Showcase Website
				</b>
				<span
					style={{
						textWrap: 'nowrap',
						textOverflow: 'ellipsis',
						overflowY: 'clip'
					}}
				>
					{props.title}
				</span>
			</div>

			<div style={{ flexGrow: 1 }}></div>

			<div
				style={{
					display: 'flex',
					alignItems: 'center',
					flexShrink: 1,
					gap: 10
				}}
				className='ms-bg-main blurry-main-background'
			>
				<NavLink
					to='/settings'
					className='ms-btn ms-rounded'
					style={{ width: 40, height: 40, position: 'relative' }}
				>
					<Settings
						size={20}
						style={{
							position: 'absolute',
							left: '50%',
							top: '50%',
							transform: 'translate(-50%, -50%)'
						}}
					/>
				</NavLink>
			</div>
		</div>
	);
}

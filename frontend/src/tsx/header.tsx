import { useState } from 'react';
import { MenuButton } from './header/menuButton';
import { HomeButton } from './header/homeButton';
import { HomeMenu } from './header/homeMenu';

export default function Header(props: { title: string }) {
	const [menuVisible, setMenuVisible] = useState(false);
	return (
		<div
			style={{
				display: 'flex',
				alignItems: 'center',
				gap: 10,
				maxWidth: '100%'
			}}
		>
			<div
				className={'ms-bg-light ms-border-main ' + (menuVisible ? 'ms-box-shadow' : '')}
				style={{
					opacity: menuVisible ? 1 : 0,
					position: 'absolute',
					zIndex: 1,
					top: 10,
					left: 10,
					overflowY: 'auto',
					maxHeight: '100%',
					height: '100%',
					margin: -10,
					paddingLeft: 9,
					paddingTop: 12,
					paddingRight: 9,
					paddingBottom: 12,
					borderWidth: 1,
					borderStyle: 'solid',
					transform: menuVisible ? 'translate(0,0)' : 'translate(-100%, 0)',
					transition: 'transform 0.3s, opacity 0.3s'
				}}
			>
				<div onClick={() => setMenuVisible(false)}>
					<div style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
						<MenuButton onClick={() => setMenuVisible(!menuVisible)} />
						<HomeButton />
					</div>
					<div style={{ marginTop: 10 }}>
						<HomeMenu />
					</div>
				</div>
			</div>
			<div style={{ display: 'flex', alignItems: 'center', gap: 5 }}>
				<MenuButton onClick={() => setMenuVisible(!menuVisible)} />
				<HomeButton />
			</div>
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
		</div>
	);
}

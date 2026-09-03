import { useState } from 'react';
import { APP_TITLE } from 'src/typescript/global';
import { PageTitle } from 'src/typescript/pageTitle';
import { HomeButton } from './header/homeButton';
import { HomeMenu } from './header/homeMenu';
import { MenuButton } from './header/menuButton';

export default function Header(props: { title: PageTitle }) {
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
					// TODO fix this layout. This div has width 0, but the intended width for it is to fill all available space of its parent
					display: 'flex',
					flexDirection: 'column',
					overflowY: 'clip',
					gap: 4,
					flexShrink: 0,
					flexBasis: 0,
					minWidth: 0,
					textWrap: 'nowrap',
					textOverflow: 'ellipsis'
				}}
			>
				<div>
					<b>{APP_TITLE}</b>
					{props.title.main ? (
						<span style={{ opacity: 0.5 }}>&nbsp;&nbsp;•&nbsp;&nbsp;</span>
					) : undefined}
					{props.title.main}
				</div>
				<div style={{ textWrap: 'nowrap', textOverflow: 'ellipsis', overflowY: 'clip' }}>
					{props.title.secondary}
				</div>
			</div>

			<div style={{ flexGrow: 1 }}></div>
		</div>
	);
}

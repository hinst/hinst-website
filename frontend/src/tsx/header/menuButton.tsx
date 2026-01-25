import { Menu } from 'react-feather';

export function MenuButton(props: { onClick: () => void }) {
	return (
		<button
			type='button'
			onClick={() => props.onClick()}
			className='ms-btn ms-primary ms-rounded'
			style={{ width: 40, height: 40, position: 'relative' }}
		>
			<Menu
				size={20}
				style={{
					position: 'absolute',
					left: '50%',
					top: '50%',
					transform: 'translate(-50%, -50%)'
				}}
			/>
		</button>
	);
}

import { NavLink } from 'react-router';
import { GoalPostRecord } from 'src/typescript/personal-goals/goalPostRecord';

interface ItemRowProps {
	item: GoalPostRecord;
}

export function ItemRow({ item }: ItemRowProps) {
	return (
		<div>
			<NavLink
				to={`/personal-goals/${item.goalId}?activePostDate=${item.dateTime}`}
				style={{ display: 'inline-flex', gap: 10 }}
			>
				<button
					type='button'
					className='ms-btn ms-primary ms-outline'
					style={{
						fontFamily: 'monospace',
						minWidth: 50,
						padding: 8,
						display: 'flex',
						justifyContent: 'center'
					}}
				>
					<span>{item.dateText}</span>
				</button>
				<div className='ms-text-main' style={{ display: 'flex', alignItems: 'center' }}>
					{item.title}
				</div>
			</NavLink>
		</div>
	);
}

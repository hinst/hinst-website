import { Calendar } from 'react-feather';
import type { ReactElement } from 'react';

export default function GoalBrowserNarrow(props: {
	activePostDate: string;
	calendarVisible: boolean;
	setCalendarVisible: (visible: boolean) => void;
	calendarTransition: string;
	getGoalCalendarPanel: () => ReactElement;
	getGoalPostPanel: () => ReactElement;
}) {
	function getFloatingCalendarButton() {
		return (
			<div
				className='ms-bg-light ms-shape-circle'
				style={{
					position: 'absolute',
					width: 40,
					height: 40,
					bottom: 0,
					right: 0,
					zIndex: 2
				}}
			>
				<button
					type='button'
					className={
						'ms-btn ms-primary ms-rounded ms-box-shadow' +
						(props.calendarVisible ? ' ms-btn-active' : '')
					}
					onClick={() => props.setCalendarVisible(!props.calendarVisible)}
					style={{
						margin: 0,
						width: 40,
						height: 40
					}}
				>
					<Calendar
						style={{
							position: 'absolute',
							left: '50%',
							top: '50%',
							transform: 'translate(-50%, -50%)'
						}}
					/>
				</button>
			</div>
		);
	}

	return (
		<div
			style={{
				position: 'relative',
				display: 'flex',
				minHeight: 0,
				height: '100%',
				maxHeight: '100%',
				width: '100%',
				maxWidth: '100%',
				overflowY: 'hidden'
			}}
		>
			{getFloatingCalendarButton()}
			<div
				className={
					'ms-bg-light ms-shape-round ms-border-main ' +
					(props.calendarVisible ? 'ms-box-shadow' : '') // Hide shadow when calendar is hidden to avoid showing ghost shadow on the right
				}
				style={{
					position: 'absolute',
					zIndex: 1,
					overflowY: 'auto',
					maxHeight: '100%',
					padding: 8,
					borderWidth: 1,
					borderStyle: 'solid',
					transform: props.calendarVisible ? 'translate(0,0)' : 'translate(-100%, 0)',
					transition: props.calendarTransition
				}}
			>
				{props.getGoalCalendarPanel()}
			</div>
			<div
				onClick={() => props.setCalendarVisible(false)}
				style={{
					display: 'flex',
					overflowY: 'auto',
					flexGrow: 1
				}}
			>
				<div style={{ width: '100%' }}>
					{props.activePostDate ? props.getGoalPostPanel() : undefined}
				</div>
			</div>
		</div>
	);
}

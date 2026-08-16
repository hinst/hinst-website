import { Calendar } from 'react-feather';
import GoalCalendarPanel from './goalCalendarPanel';
import GoalPostPanel from './goalPostPanel';
import { GoalPostHeaderEx } from 'src/typescript/rest_objects/goalPostHeaderEx';

export default function GoalBrowserNarrow(props: {
	goalId: string;
	activePostDate: string;
	calendarVisible: boolean;
	setCalendarVisible: (visible: boolean) => void;
	calendarTransition: string;
	goalManagerMode: boolean;
	receivePosts: (posts: GoalPostHeaderEx[]) => void;
	reloadGoalCalendar: number;
	onChangeGoalCalendar: () => void;
}) {
	const {
		goalId,
		activePostDate,
		calendarVisible,
		setCalendarVisible,
		calendarTransition,
		goalManagerMode,
		receivePosts,
		reloadGoalCalendar,
		onChangeGoalCalendar
	} = props;

	function getGoalCalendarPanel() {
		return (
			<GoalCalendarPanel
				id={goalId}
				receivePosts={receivePosts}
				activePostDate={parseInt(activePostDate) || 0}
				reload={reloadGoalCalendar}
			/>
		);
	}

	function getGoalPostPanel() {
		return (
			<GoalPostPanel
				goalId={parseInt(goalId)}
				postDate={parseInt(activePostDate)}
				goalManagerMode={goalManagerMode}
				onChange={onChangeGoalCalendar}
			/>
		);
	}

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
						(calendarVisible ? ' ms-btn-active' : '')
					}
					onClick={() => setCalendarVisible(!calendarVisible)}
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
					(calendarVisible ? 'ms-box-shadow' : '') // Hide shadow when calendar is hidden to avoid showing ghost shadow on the right
				}
				style={{
					position: 'absolute',
					zIndex: 1,
					overflowY: 'auto',
					maxHeight: '100%',
					padding: 8,
					borderWidth: 1,
					borderStyle: 'solid',
					transform: calendarVisible ? 'translate(0,0)' : 'translate(-100%, 0)',
					transition: calendarTransition
				}}
			>
				{getGoalCalendarPanel()}
			</div>
			<div
				onClick={() => setCalendarVisible(false)}
				style={{
					display: 'flex',
					overflowY: 'auto',
					flexGrow: 1
				}}
			>
				<div style={{ width: '100%' }}>
					{activePostDate ? getGoalPostPanel() : undefined}
				</div>
			</div>
		</div>
	);
}

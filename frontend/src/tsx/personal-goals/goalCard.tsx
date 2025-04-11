import { NavLink } from 'react-router';
import { GoalRecord } from 'src/typescript/personal-goals/goalHeader';
import { GOAL_INFOS, translateGoalTitle } from 'src/typescript/personal-goals/goalInfo';
import { getCurrentLanguage } from 'src/typescript/language';

export function GoalCard(props: {goal: GoalRecord}) {
	const currentLanguage = getCurrentLanguage();
	return <div
		style={{
			width: 'fit-content',
			padding: 0,
			margin: 0,
			overflow: 'hidden'
		}}
		className='ms-card ms-border grayscale'
	>
		<NavLink
			to={'/personal-goals/' + props.goal.id}
			key={props.goal.id}
			className='ms-text-main'
		>
			<div
				style={{display: 'flex', flexDirection: 'column'}}
				className='ms-bg-light'
			>
				<img
					width={200}
					height={100}
					src={GOAL_INFOS.get(props.goal.title)?.coverImage}
					alt={props.goal.title}
				/>
				<div
					style={{
						padding: 8,
						maxWidth: 200,
					}}
				>
					{translateGoalTitle(currentLanguage, props.goal.title)}
				</div>
			</div>
		</NavLink>
	</div>
}
import { NavLink } from 'react-router';
import { GoalHeader } from './goalHeader';
import { GOAL_INFOS, translateGoalTitle } from './goalInfo';

export default function GoalList(props: {goals: GoalHeader[]}) {
	return <div>
		{props.goals.map(goal =>
			<div
				key={goal.id}
				style={{
					width: 'fit-content',
					padding: 0,
					overflow: 'hidden'
				}}
				className='ms-card ms-border grayscale'
			>
				<NavLink
					to={'/personal-goals/' + goal.id}
					key={goal.id}
					className='ms-text-main'
				>
					<div
						style={{display: 'flex', flexDirection: 'column'}}
						className='ms-bg-light'
					>
						<img
							width={200}
							height={100}
							src={GOAL_INFOS.get(goal.title)?.coverImage}
							alt={goal.title}
						/>
						<div
							style={{margin: 8}}
						>
							{translateGoalTitle(goal.title)}
						</div>
					</div>
				</NavLink>
			</div>
		)}
	</div>;
}
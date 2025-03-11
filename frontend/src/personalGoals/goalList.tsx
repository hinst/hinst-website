import { NavLink } from 'react-router';
import { GoalHeader } from './goalHeader';
import { translateGoalTitle } from './goalTitle';

export default function GoalList(props: {goals: GoalHeader[]}) {
	return <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
		{props.goals.map(goal =>
			<div key={goal.id}>
				<NavLink className='ms-btn ms-primary ms-outline' to={'/personal-goals/' + goal.id} key={goal.id}>
					{translateGoalTitle(goal.title)}
				</NavLink>
			</div>
		)}
	</div>;
}
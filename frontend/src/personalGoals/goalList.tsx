import { NavLink } from 'react-router';
import { GoalHeader } from './goalHeader';

export default function GoalList(props: {goals: GoalHeader[]}) {
	return <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
		{props.goals.map(goal =>
			<div>
				<NavLink className='ms-btn ms-primary ms-outline' to={'/personal-goals/' + goal.id} key={goal.id}>
					{goal.title}
				</NavLink>
			</div>
		)}
	</div>;
}
import GoalListPanel from './personalGoals/goalListPanel';

export default function HomePage() {
	// <NavLink to='/personal-goals' className='ms-btn ms-outline ms-primary'>My Personal Goals</NavLink>
	return <div>
		<GoalListPanel />
	</div>;
}
import { GoalObjectWithMethods } from 'src/typescript/rest_objects/restObjectExtensions';
import { GoalCard } from './goalCard';

export default function GoalList(props: { goals: GoalObjectWithMethods[] }) {
	return (
		<div style={{ display: 'flex', flexWrap: 'wrap', gap: 10, flexDirection: 'row' }}>
			{props.goals.map((goal) => (
				<GoalCard key={goal.id} goal={goal} />
			))}
		</div>
	);
}

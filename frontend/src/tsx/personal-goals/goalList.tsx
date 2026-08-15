import { GoalObject } from 'src/typescript/generated/rest_objects';
import 'src/typescript/restObjectExtensions';
import { GoalCard } from './goalCard';

export default function GoalList(props: { goals: GoalObject[] }) {
	return (
		<div style={{ display: 'flex', flexWrap: 'wrap', gap: 10, flexDirection: 'row' }}>
			{props.goals.map((goal) => (
				<GoalCard key={goal.id} goal={goal} />
			))}
		</div>
	);
}

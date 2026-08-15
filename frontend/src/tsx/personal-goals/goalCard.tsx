import { NavLink } from 'react-router';
import { GoalObjectEx } from 'src/typescript/rest_objects/restObjectExtensions';
import { AppContext } from '../context';
import { useContext } from 'react';
import { apiClient } from 'src/typescript/apiClient';

export function GoalCard(props: { goal: GoalObjectEx }) {
	const context = useContext(AppContext);
	return (
		<div
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
				<div style={{ display: 'flex', flexDirection: 'column' }} className='ms-bg-light'>
					<img
						width={200}
						height={100}
						src={apiClient.getGoalImageUrl(props.goal.id)}
						alt={props.goal.title}
					/>
					<div
						style={{
							padding: 8,
							maxWidth: 200
						}}
					>
						{props.goal.getTitle(context.currentLanguage)}
					</div>
				</div>
			</NavLink>
		</div>
	);
}

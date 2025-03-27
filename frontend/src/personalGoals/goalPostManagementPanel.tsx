import { Tool } from 'react-feather';
import { SmartPostExtended } from './smartPost';
import { API_URL } from '../global';

export default function GoalPostManagementPanel(props: {
	postData: SmartPostExtended,
	setPostData: (postData: SmartPostExtended) => void,
}) {
	async function setPublic(isPublic: boolean) {
		const url = API_URL + '/goalPost/setPublic' +
			'?goalId=' + encodeURIComponent(props.postData.obj_id) +
			'&postDateTime=' + encodeURIComponent(props.postData.date) +
			'&isPublic=' + encodeURIComponent('' + isPublic);
		const response = await fetch(url);
		if (!response.ok)
			throw new Error('Cannot update post visibility. Status: ' + response.statusText);
		props.setPostData({...props.postData, isPublic});
	}

	return <div
		className='ms-alert ms-light'
		style={{
			display: 'flex',
			gap: 10,
			alignItems: 'center',
		}}
	>
		<Tool/>
		<input
			type='checkbox'
			checked={props.postData?.isPublic}
			onChange={() => setPublic(!props.postData?.isPublic)}
		/>
		public
	</div>;
}
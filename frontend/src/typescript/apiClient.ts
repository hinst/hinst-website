import { GoalPostObject } from './personal-goals/smartPost';

class ApiClient {
	url: string = process.env.API_URL || '/hinst-website/api';

	async fetch(url: string, options?: RequestInit): Promise<Response> {
		const response = await fetch(this.url + url, options);
		return response;
	}

	async goalPostSetPublic(goalId: number, postDateTime: number, isPublic: boolean): Promise<Response> {
		const url = '/goalPost/setPublic' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime) +
			'&isPublic=' +
			encodeURIComponent('' + isPublic);
		return this.fetch(url);
	}

	async getGoalPost(goalId: number, postDateTime: number): Promise<GoalPostObject> {
		const url = '/goalPost' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime);
		const response = await this.fetch(url);
		if (!response.ok)
			throw new Error('Cannot load post. Status: ' + response.statusText);
		return await response.json() as GoalPostObject;
	}
}

export const apiClient = new ApiClient();
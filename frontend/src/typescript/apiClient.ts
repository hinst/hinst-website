import {
	GoalObject,
	GoalPostHeader,
	GoalPostObject,
	GoalPostSearchIndexingHeader
} from 'src/typescript/generated/rest_objects';
import {
	goalObjectWithMethods,
	goalPostHeaderWithMethods,
	GoalObjectWithMethods,
	GoalPostHeaderWithMethods
} from './rest_objects/restObjectExtensions';
import { RiddleItem } from './riddle';
import { settingsStorage } from './settings';

class ApiClient {
	readonly url: string = process.env.API_URL || '/hinst-website/api';

	private async fetch(url: string, options?: RequestInit): Promise<Response> {
		if (settingsStorage.language) {
			options = options || {};
			options.headers = {
				...options.headers,
				'Accept-Language': settingsStorage.language
			};
		}
		const response = await fetchSafe(this.url + url, options);
		return response;
	}

	async getGoal(id: number): Promise<GoalObjectWithMethods> {
		const url = '/goal?id=' + encodeURIComponent(id);
		const response = await this.fetch(url);
		const data = (await response.json()) as GoalObject;
		return goalObjectWithMethods(data);
	}

	async getGoals(): Promise<GoalObjectWithMethods[]> {
		const response = await this.fetch('/goals');
		const data = ((await response.json()) as GoalObject[]) || [];
		return data.map((item) => goalObjectWithMethods(item));
	}

	async goalPostSetPublic(
		goalId: number,
		postDateTime: number,
		isPublic: boolean
	): Promise<Response> {
		const url =
			'/goalPost/setPublic' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime) +
			'&isPublic=' +
			encodeURIComponent('' + isPublic);
		return await this.fetch(url);
	}

	async goalPostSetSearchIndexingEnabled(
		goalId: number,
		postDateTime: number,
		enabled: boolean
	): Promise<Response> {
		const url =
			'/goalPost/setSearchIndexingEnabled' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime) +
			'&enabled=' +
			encodeURIComponent('' + enabled);
		return await this.fetch(url);
	}

	async getGoalPost(goalId: number, postDateTime: number): Promise<GoalPostObject> {
		const url =
			'/goalPost' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime);
		const response = await this.fetch(url);
		return (await response.json()) as GoalPostObject;
	}

	async setGoalPostText(
		goalId: number,
		postDateTime: number,
		languageTag: string,
		text: string
	): Promise<Response> {
		const url =
			'/goalPost/setText' +
			'?goalId=' +
			encodeURIComponent(goalId) +
			'&postDateTime=' +
			encodeURIComponent(postDateTime) +
			'&languageTag=' +
			encodeURIComponent(languageTag);
		return this.fetch(url, { method: 'POST', body: text });
	}

	async createRiddle(): Promise<RiddleItem> {
		const url = '/riddles/new';
		const response = await this.fetch(url);
		const object = await response.json();
		return Object.assign(new RiddleItem(), object);
	}

	async getPrimeNumbers(): Promise<number[]> {
		const url = '/riddles/primeNumbers';
		const response = await this.fetch(url);
		return (await response.json()) as number[];
	}

	getImageUrl(goalId: number, postDateTime: number, index: number): string {
		const goalIdParameter = encodeURIComponent('' + goalId);
		const postDateTimeParameter = encodeURIComponent('' + postDateTime);
		const indexParameter = encodeURIComponent('' + index);
		return (
			this.url +
			'/goalPost/image?goalId=' +
			goalIdParameter +
			'&postDateTime=' +
			postDateTimeParameter +
			'&index=' +
			indexParameter
		);
	}

	getGoalImageUrl(goalId: number): string {
		return this.url + '/goal/image?id=' + goalId;
	}

	async getUrlPings(): Promise<GoalPostSearchIndexingHeader[]> {
		const url = '/urlPings';
		const response = await this.fetch(url);
		return ((await response.json()) as GoalPostSearchIndexingHeader[]) || [];
	}

	async pingUrlManually(url: string) {
		const apiUrl = '/pingUrlManually' + '?url=' + encodeURIComponent(url);
		return await this.fetch(apiUrl, { method: 'PUT' });
	}

	async getGoalPosts(goalId: number): Promise<GoalPostHeaderWithMethods[]> {
		const url = '/goalPosts?id=' + encodeURIComponent(goalId);
		const response = await this.fetch(url);
		const items = ((await response.json()) as GoalPostHeader[]) || [];
		return items
			.filter((item) => item.type === 'post')
			.map((item) => goalPostHeaderWithMethods(item));
	}

	async searchGoalPosts(query: string): Promise<GoalPostHeaderWithMethods[]> {
		const url = '/goalPosts/search?query=' + encodeURIComponent(query);
		const response = await this.fetch(url);
		const items = ((await response.json()) as GoalPostHeader[]) || [];
		return items.map((item) => goalPostHeaderWithMethods(item));
	}
}

export const apiClient = new ApiClient();

async function fetchSafe(url: string, requestInit: RequestInit = {}): Promise<Response> {
	const response = await fetch(url, requestInit);
	if (!response.ok) throw new Error(url + ' ' + response.statusText);
	return response;
}

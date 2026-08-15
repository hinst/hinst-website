import {
	GoalObject,
	GoalPostHeader,
	GoalPostObject,
	GoalPostSearchIndexingHeader
} from 'src/typescript/generated/rest_objects';
import { goalObjectWithMethods, GoalObjectEx } from './rest_objects/goalObjectEx';
import { goalPostHeaderWithMethods, GoalPostHeaderEx } from './rest_objects/goalPostHeaderEx';
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

	async getGoal(id: number): Promise<GoalObjectEx> {
		const url = '/goal?' + new URLSearchParams({ id: '' + id });
		const response = await this.fetch(url);
		const data = (await response.json()) as GoalObject;
		return goalObjectWithMethods(data);
	}

	async getGoals(): Promise<GoalObjectEx[]> {
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
			'/goalPost/setPublic?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime,
				isPublic: '' + isPublic
			});
		return await this.fetch(url);
	}

	async goalPostSetSearchIndexingEnabled(
		goalId: number,
		postDateTime: number,
		enabled: boolean
	): Promise<Response> {
		const url =
			'/goalPost/setSearchIndexingEnabled?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime,
				enabled: '' + enabled
			});
		return await this.fetch(url);
	}

	async getGoalPost(goalId: number, postDateTime: number): Promise<GoalPostObject> {
		const url =
			'/goalPost?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime
			});
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
			'/goalPost/setText?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime,
				languageTag
			});
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
		return (
			this.url +
			'/goalPost/image?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime,
				index: '' + index
			})
		);
	}

	getGoalImageUrl(goalId: number): string {
		return this.url + '/goal/image?' + new URLSearchParams({ id: '' + goalId });
	}

	async getUrlPings(): Promise<GoalPostSearchIndexingHeader[]> {
		const url = '/urlPings';
		const response = await this.fetch(url);
		return ((await response.json()) as GoalPostSearchIndexingHeader[]) || [];
	}

	async pingUrlManually(url: string) {
		const apiUrl = '/pingUrlManually?' + new URLSearchParams({ url });
		return await this.fetch(apiUrl, { method: 'PUT' });
	}

	async getGoalPosts(goalId: number): Promise<GoalPostHeaderEx[]> {
		const url = '/goalPosts?' + new URLSearchParams({ id: '' + goalId });
		const response = await this.fetch(url);
		const items = ((await response.json()) as GoalPostHeader[]) || [];
		return items
			.filter((item) => item.type === 'post')
			.map((item) => goalPostHeaderWithMethods(item));
	}

	async searchGoalPosts(query: string): Promise<GoalPostHeaderEx[]> {
		const url = '/goalPosts/search?' + new URLSearchParams({ query });
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

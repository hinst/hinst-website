import { GoalPostObject } from './personal-goals/smartPost';
import { RiddleItem } from './riddle';
import { settingsStorage } from './settings';
import { UrlPingRecord } from './urlPing';

class ApiClient {
	readonly url: string = process.env.API_URL || '/hinst-website/api';

	async fetch(url: string, options?: RequestInit): Promise<Response> {
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

	async getUrlPings(): Promise<UrlPingRecord[]> {
		const url = '/urlPings';
		const response = await this.fetch(url);
		return (await response.json()) as UrlPingRecord[];
	}
}

export const apiClient = new ApiClient();

async function fetchSafe(url: string, requestInit: RequestInit = {}): Promise<Response> {
	const response = await fetch(url, requestInit);
	if (!response.ok) throw new Error(url + ' ' + response.statusText);
	return response;
}

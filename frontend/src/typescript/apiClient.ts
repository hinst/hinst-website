import createClient from 'openapi-fetch';
import type { paths } from 'src/typescript/generated/openapi';
import type { GoalPostObject, GoalPostSearchIndexingHeader } from './apiTypes';
import { createGoalObjectEx, type GoalObjectEx } from './rest_objects/goalObjectEx';
import { createGoalPostHeaderEx, type GoalPostHeaderEx } from './rest_objects/goalPostHeaderEx';
import { settingsStorage } from './settings';

const API_URL = '/hinst-website/api';
const apiUrl = process.env.API_URL || API_URL;
const baseApiUrl = apiUrl.endsWith(API_URL)
	? apiUrl.substring(0, apiUrl.length - API_URL.length)
	: apiUrl;

const client = createClient<paths>({
	baseUrl: baseApiUrl,
	bodySerializer: (body) => (typeof body === 'string' ? body : JSON.stringify(body)),
});

client.use({
	onRequest: (params) => {
		const language = settingsStorage.language;
		if (language) params.request.headers.set('Accept-Language', language);
	},
	onResponse: (params) => {
		if (!params.response.ok)
			throw new Error(
				params.request.url + ' ' + params.response.status + ' ' + params.response.statusText,
			);
	},
});

function requireData<T>(data: T | undefined | null): T {
	if (data === null || data === undefined) throw new Error('API client: data is missing');
	return data;
}

class ApiClient {
	readonly url: string = apiUrl;

	async getGoal(id: number): Promise<GoalObjectEx> {
		const { data } = await client.GET('/hinst-website/api/goal', {
			params: { query: { id } },
		});
		return createGoalObjectEx(requireData(data));
	}

	async getGoals(): Promise<GoalObjectEx[]> {
		const { data } = await client.GET('/hinst-website/api/goals');
		return requireData(data).map((goal) => createGoalObjectEx(goal));
	}

	async goalPostSetPublic(
		goalId: number,
		postDateTime: number,
		isPublic: boolean,
	): Promise<Response> {
		const { response } = await client.PUT('/hinst-website/api/goalPost/setPublic', {
			params: { query: { goalId, postDateTime, isPublic } },
		});
		return response;
	}

	async goalPostSetSearchIndexingEnabled(
		goalId: number,
		postDateTime: number,
		enabled: boolean,
	): Promise<Response> {
		const { response } = await client.PUT('/hinst-website/api/goalPost/setSearchIndexingEnabled', {
			params: { query: { goalId, postDateTime, enabled } },
		});
		return response;
	}

	async getGoalPost(goalId: number, postDateTime: number): Promise<GoalPostObject> {
		const { data } = await client.GET('/hinst-website/api/goalPost', {
			params: { query: { goalId, postDateTime } },
		});
		return requireData(data);
	}

	async setGoalPostText(
		goalId: number,
		postDateTime: number,
		languageTag: string,
		text: string,
	): Promise<Response> {
		const { response } = await client.POST('/hinst-website/api/goalPost/setText', {
			params: { query: { goalId, postDateTime, languageTag } },
			body: text,
			headers: { 'Content-Type': 'text/plain;charset=UTF-8' },
		});
		return response;
	}

	getImageUrl(goalId: number, postDateTime: number, index: number): string {
		return (
			this.url +
			'/goalPost/image?' +
			new URLSearchParams({
				goalId: '' + goalId,
				postDateTime: '' + postDateTime,
				index: '' + index,
			})
		);
	}

	getGoalImageUrl(goalId: number): string {
		return this.url + '/goal/image?' + new URLSearchParams({ id: '' + goalId });
	}

	async getUrlPings(): Promise<GoalPostSearchIndexingHeader[]> {
		const { data } = await client.GET('/hinst-website/api/urlPings');
		return requireData(data);
	}

	async getGoalPosts(goalId: number): Promise<GoalPostHeaderEx[]> {
		const { data } = await client.GET('/hinst-website/api/goalPosts', {
			params: { query: { id: goalId } },
		});
		return requireData(data)
			.filter((post) => post.type === 'post')
			.map((post) => createGoalPostHeaderEx(post));
	}

	async searchGoalPosts(query: string): Promise<GoalPostHeaderEx[]> {
		const { data } = await client.GET('/hinst-website/api/goalPosts/search', {
			params: { query: { query } },
		});
		if (!data) throw new Error('Missing data');
		return data.map((post) => createGoalPostHeaderEx(post));
	}

	async isAdminModeEnabled(): Promise<boolean> {
		const { data } = await client.GET('/hinst-website/api/isAdminModeEnabled');
		return requireData(data);
	}

	async setGoalPostGooglePingedAt(goalId: number, postDateTime: number): Promise<boolean> {
		const { data } = await client.PUT('/hinst-website/api/goalPosts/googlePingedAt', {
			params: { query: { goalId, postDateTime } },
		});
		return requireData(data);
	}
}

export const apiClient = new ApiClient();

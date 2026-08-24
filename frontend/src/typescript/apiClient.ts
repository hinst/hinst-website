import createClient from 'openapi-fetch';
import type { paths } from 'src/typescript/generated/openapi';
import { goalObjectWithMethods, GoalObjectEx } from './rest_objects/goalObjectEx';
import { goalPostHeaderWithMethods, GoalPostHeaderEx } from './rest_objects/goalPostHeaderEx';
import type {
	GoalObject,
	GoalPostHeader,
	GoalPostObject,
	GoalPostSearchIndexingHeader
} from './apiTypes';
import { settingsStorage } from './settings';

/** Base URL of the REST API, e.g. '/hinst-website/api' or 'https://host/hinst-website/api'. */
const apiRoot = process.env.API_URL || '/hinst-website/api';

// The paths in the OpenAPI spec include the API root (e.g. '/hinst-website/api/goal'),
// so the client base URL is the part of the API root before that.
const client = createClient<paths>({
	baseUrl: apiRoot.replace(/\/hinst-website\/api\/?$/, ''),
	// Send string bodies as-is (e.g. goal post text), everything else as JSON
	bodySerializer: (body) => (typeof body === 'string' ? body : JSON.stringify(body))
});

// Add the Accept-Language header to every request and throw on non-2xx responses
client.use({
	onRequest: (params) => {
		const language = settingsStorage.language;
		if (language) params.request.headers.set('Accept-Language', language);
	},
	onResponse: (params) => {
		if (!params.response.ok)
			throw new Error(
				params.request.url + ' ' + params.response.status + ' ' + params.response.statusText
			);
	}
});

class ApiClient {
	readonly url: string = apiRoot;

	async getGoal(id: number): Promise<GoalObjectEx> {
		const { data } = await client.GET('/hinst-website/api/goal', {
			params: { query: { id } }
		});
		return goalObjectWithMethods(data!);
	}

	async getGoals(): Promise<GoalObjectEx[]> {
		const { data } = await client.GET('/hinst-website/api/goals');
		const goals = data!;
		return goals.map((goal) => goalObjectWithMethods(goal));
	}

	async goalPostSetPublic(
		goalId: number,
		postDateTime: number,
		isPublic: boolean
	): Promise<Response> {
		const { response } = await client.PUT('/hinst-website/api/goalPost/setPublic', {
			params: { query: { goalId, postDateTime, isPublic } }
		});
		return response;
	}

	async goalPostSetSearchIndexingEnabled(
		goalId: number,
		postDateTime: number,
		enabled: boolean
	): Promise<Response> {
		const { response } = await client.PUT('/hinst-website/api/goalPost/setSearchIndexingEnabled', {
			params: { query: { goalId, postDateTime, enabled } }
		});
		return response;
	}

	async getGoalPost(goalId: number, postDateTime: number): Promise<GoalPostObject> {
		const { data } = await client.GET('/hinst-website/api/goalPost', {
			params: { query: { goalId, postDateTime } }
		});
		return data!;
	}

	async setGoalPostText(
		goalId: number,
		postDateTime: number,
		languageTag: string,
		text: string
	): Promise<Response> {
		const { response } = await client.POST('/hinst-website/api/goalPost/setText', {
			params: { query: { goalId, postDateTime, languageTag } },
			body: text,
			headers: { 'Content-Type': 'text/plain;charset=UTF-8' }
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
				index: '' + index
			})
		);
	}

	getGoalImageUrl(goalId: number): string {
		return this.url + '/goal/image?' + new URLSearchParams({ id: '' + goalId });
	}

	async getUrlPings(): Promise<GoalPostSearchIndexingHeader[]> {
		const { data } = await client.GET('/hinst-website/api/urlPings');
		return data!;
	}

	async getGoalPosts(goalId: number): Promise<GoalPostHeaderEx[]> {
		const { data } = await client.GET('/hinst-website/api/goalPosts', {
			params: { query: { id: goalId } }
		});
		const posts = data!;
		return posts
			.filter((post) => post.type === 'post')
			.map((post) => goalPostHeaderWithMethods(post));
	}

	async searchGoalPosts(query: string): Promise<GoalPostHeaderEx[]> {
		const { data } = await client.GET('/hinst-website/api/goalPosts/search', {
			params: { query: { query } }
		});
		const posts = data!;
		return posts.map((post) => goalPostHeaderWithMethods(post));
	}
}

export const apiClient = new ApiClient();

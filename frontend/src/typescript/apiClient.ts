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
}

export const apiClient = new ApiClient();
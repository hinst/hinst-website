//@ts-check
const BASE_URL = 'http://192.168.0.23/hinst-website';
const GOAL_ID = 247488;

/** @type {Record<string, string>} */
const defaultHeaders = {
	'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
	'cache-control': 'no-cache',
	pragma: 'no-cache',
	Referer: `${BASE_URL}/`
};

export class Requests {
	/** @type {number[]} */
	dateTimes = [];
	/** @type {number} */
	lastImageCount = 0;
	/** @type {number} */
	requestCount = 0;

	/**
		@param {string} path
		@param {string} accept
		@param {string} label
	*/
	async fetchResource(path, accept, label) {
		this.requestCount++;
		const response = await fetch(`${BASE_URL}${path}`, {
			headers: { ...defaultHeaders, accept },
			method: 'GET'
		});
		if (response.status !== 200) throw new Error(`${label} returned status ${response.status}`);
		const buffer = await response.arrayBuffer();
		return buffer.byteLength;
	}

	/**
		@param {string} path
		@param {string} accept
		@param {string} label
		@param {object} [extraHeaders]
		@returns {Promise<any>}
	*/
	async fetchJSON(path, accept, label, extraHeaders = {}) {
		this.requestCount++;
		const response = await fetch(`${BASE_URL}${path}`, {
			headers: { ...defaultHeaders, accept, ...extraHeaders },
			method: 'GET'
		});
		if (response.status !== 200) throw new Error(`${label} returned status ${response.status}`);
		return response.json();
	}

	async main() {
		return this.fetchResource(
			'/',
			'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
			'main'
		);
	}

	async css1() {
		return this.fetchResource('/app.3e369f90.css', 'text/css,*/*;q=0.1', 'css1');
	}

	async css2() {
		return this.fetchResource('/app.bef21471.css', 'text/css,*/*;q=0.1', 'css2');
	}

	async javaScript() {
		return this.fetchResource('/app.9da11e27.js', '*/*', 'javaScript');
	}

	async icon() {
		return this.fetchResource(
			'/icon.fd8aa8a2.webp',
			'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
			'icon'
		);
	}

	async getGoalPosts() {
		const data = await this.fetchJSON(`/api/goalPosts?id=${GOAL_ID}`, '*/*', 'getGoalPosts');
		if (!Array.isArray(data))
			throw new Error(`getGoalPosts: expected array, got ${typeof data}`);
		this.dateTimes = data.map((item) => item.dateTime);
		return new TextEncoder().encode(JSON.stringify(data)).byteLength;
	}

	async getGoal() {
		return this.fetchResource(`/api/goal?id=${GOAL_ID}`, '*/*', 'getGoal');
	}

	/**
		@param {number} postDateTime
	*/
	async getGoalPost(postDateTime) {
		const data = await this.fetchJSON(
			`/api/goalPost?goalId=${GOAL_ID}&postDateTime=${postDateTime}`,
			'*/*',
			'getGoalPost',
			{ 'accept-language': 'ru' }
		);
		this.lastImageCount = data.imageCount;
		return new TextEncoder().encode(JSON.stringify(data)).byteLength;
	}

	async favicon() {
		return this.fetchResource(
			'/favicon.ce92f6f7.png',
			'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
			'favicon'
		);
	}

	/**
		@param {number} postDateTime
		@param {number} index
	*/
	async image(postDateTime, index) {
		return this.fetchResource(
			`/api/goalPost/image?goalId=${GOAL_ID}&postDateTime=${postDateTime}&index=${index}`,
			'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
			'image'
		);
	}

	async all() {
		const methods = ['main', 'css1', 'css2', 'javaScript', 'icon', 'getGoalPosts', 'getGoal'];
		const sizes = [];
		for (const name of methods) {
			sizes.push(await this[name]());
		}
		for (const dateTime of this.dateTimes) {
			sizes.push(await this.getGoalPost(dateTime));

			for (let i = 0; i < this.lastImageCount; i++) {
				sizes.push(await this.image(dateTime, i));
			}
		}
		sizes.push(await this.favicon());

		return sizes;
	}
}

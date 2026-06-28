//@ts-check
export class Requests {
	/** @type {number[]} */
	dateTimes = [];
	/** @type {number} */
	lastImageCount = 0;

	async main() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/', {
			headers: {
				accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				'upgrade-insecure-requests': '1'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`main returned status ${response.status}`);
		return buffer.byteLength;
	}

	async css1() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.3e369f90.css', {
			headers: {
				accept: 'text/css,*/*;q=0.1',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`css1 returned status ${response.status}`);
		return buffer.byteLength;
	}

	async css2() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.bef21471.css', {
			headers: {
				accept: 'text/css,*/*;q=0.1',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`css2 returned status ${response.status}`);
		return buffer.byteLength;
	}

	async javaScript() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.9da11e27.js', {
			headers: {
				accept: '*/*',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`javaScript returned status ${response.status}`);
		return buffer.byteLength;
	}

	async icon() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/icon.fd8aa8a2.webp', {
			headers: {
				accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`icon returned status ${response.status}`);
		return buffer.byteLength;
	}

	async api1() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/api/goalPosts?id=247488',
			{
				headers: {
					accept: '*/*',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				method: 'GET'
			}
		);
		const text = await response.text();
		if (response.status !== 200) throw new Error(`api1 returned status ${response.status}`);
		const data = JSON.parse(text);
		if (!Array.isArray(data)) throw new Error(`api1: expected array, got ${typeof data}`);
		this.dateTimes = data.map((item) => item.dateTime);
		return new TextEncoder().encode(text).byteLength;
	}

	async api2() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/api/goal?id=247488', {
			headers: {
				accept: '*/*',
				'accept-language': 'ru',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			method: 'GET'
		});
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`api2 returned status ${response.status}`);
		return buffer.byteLength;
	}

	/**
		@param {number} postDateTime
	*/
	async api3(postDateTime) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost?goalId=247488&postDateTime=${postDateTime}`,
			{
				headers: {
					accept: '*/*',
					'accept-language': 'ru',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				method: 'GET'
			}
		);
		const text = await response.text();
		if (response.status !== 200) throw new Error(`api3 returned status ${response.status}`);
		const data = JSON.parse(text);
		this.lastImageCount = data.imageCount;
		return new TextEncoder().encode(text).byteLength;
	}

	async favicon() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/favicon.ce92f6f7.png',
			{
				headers: {
					accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				method: 'GET'
			}
		);
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`favicon returned status ${response.status}`);
		return buffer.byteLength;
	}

	/**
		@param {number} postDateTime
		@param {number} index
	*/
	async image(postDateTime, index) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=${postDateTime}&index=${index}`,
			{
				headers: {
					accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				method: 'GET'
			}
		);
		const buffer = await response.arrayBuffer();
		if (response.status !== 200) throw new Error(`image returned status ${response.status}`);
		return buffer.byteLength;
	}

	async all() {
		const methods = ['main', 'css1', 'css2', 'javaScript', 'icon', 'api1', 'api2'];
		const sizes = [];
		for (const name of methods) {
			sizes.push(await this[name]());
		}
		for (const dateTime of this.dateTimes) {
			sizes.push(await this.api3(dateTime));

			for (let i = 0; i < this.lastImageCount; i++) {
				sizes.push(await this.image(dateTime, i));
			}
		}
		sizes.push(await this.favicon());

		return sizes;
	}
}

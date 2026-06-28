//@ts-check
export class Requests {
	/** @type {number[]} */
	dateTimes = [];

	async main() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/', {
			headers: {
				accept: 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				'upgrade-insecure-requests': '1',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.text()];
	}

	async css1() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.3e369f90.css', {
			headers: {
				accept: 'text/css,*/*;q=0.1',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.text()];
	}

	async css2() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.bef21471.css', {
			headers: {
				accept: 'text/css,*/*;q=0.1',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.text()];
	}

	async javaScript() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/app.9da11e27.js', {
			headers: {
				accept: '*/*',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.text()];
	}

	async icon() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/icon.fd8aa8a2.webp', {
			headers: {
				accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
				'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.arrayBuffer()];
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
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		const data = await response.json();
		this.dateTimes = data.map(item => item.dateTime);
		return [response.status, data];
	}

	async api2() {
		const response = await fetch('http://192.168.0.23:30001/hinst-website/api/goal?id=247488', {
			headers: {
				accept: '*/*',
				'accept-language': 'ru',
				'cache-control': 'no-cache',
				pragma: 'no-cache',
				cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
				Referer: 'http://192.168.0.23:30001/hinst-website/'
			},
			body: null,
			method: 'GET'
		});
		return [response.status, await response.json()];
	}

	async api3(postDateTime) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost?goalId=247488&postDateTime=${postDateTime}`,
			{
				headers: {
					accept: '*/*',
					'accept-language': 'ru',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		return [response.status, await response.json()];
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
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		return [response.status, await response.arrayBuffer()];
	}

	/**
		@param {number} postDateTime
	*/
	async image1(postDateTime) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=${postDateTime}&index=0`,
			{
				headers: {
					accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		return [response.status, await response.arrayBuffer()];
	}

	/**
		@param {number} postDateTime
	*/
	async image2(postDateTime) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=${postDateTime}&index=1`,
			{
				headers: {
					accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		return [response.status, await response.arrayBuffer()];
	}

	/**
		@param {number} postDateTime
	*/
	async image3(postDateTime) {
		const response = await fetch(
			`http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=${postDateTime}&index=2`,
			{
				headers: {
					accept: 'image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8',
					'accept-language': 'en,ru;q=0.9,en-GB;q=0.8,en-US;q=0.7,de;q=0.6',
					'cache-control': 'no-cache',
					pragma: 'no-cache',
					cookie: 'CSRF-Token-AYTHB54=CJF9aQ4o4abMGtjmsSMfwRfzA6TtUts7F2H5XhZCrMvmwRNayveChbaGUE5tr2da; sessionid-AYTHB54=rUuuDpYYT4utPG7C2eREsXaGVEttCNFh2She3kps36U6Y5Rv7zgFmDVs7mQgLNYg',
					Referer: 'http://192.168.0.23:30001/hinst-website/'
				},
				body: null,
				method: 'GET'
			}
		);
		return [response.status, await response.arrayBuffer()];
	}

	async all() {
		const methods = [
			'main',
			'css1',
			'css2',
			'javaScript',
			'icon',
			'api1',
			'api2'
		];
		const sizes = [];
		for (const name of methods) {
			const [status, response] = await this[name]();
			let responseSize = -1;
			if (typeof response === 'string') responseSize = response.length;
			else if (response instanceof ArrayBuffer) responseSize = response.byteLength;
			else if (typeof response === 'object') responseSize = JSON.stringify(response).length;
			if (status === 200) sizes.push(responseSize);
			else throw new Error(`${name} returned status ${status}`);
		}
		for (const dateTime of this.dateTimes) {
			const [status, response] = await this.api3(dateTime);
			let responseSize = -1;
			if (typeof response === 'string') responseSize = response.length;
			else if (response instanceof ArrayBuffer) responseSize = response.byteLength;
			else if (typeof response === 'object') responseSize = JSON.stringify(response).length;
			if (status === 200) sizes.push(responseSize);
			else throw new Error(`api3 returned status ${status}`);

			for (const imageMethod of ['image1', 'image2', 'image3']) {
				const [imgStatus, imgResponse] = await this[imageMethod](dateTime);
				let imgSize = -1;
				if (typeof imgResponse === 'string') imgSize = imgResponse.length;
				else if (imgResponse instanceof ArrayBuffer) imgSize = imgResponse.byteLength;
				else if (typeof imgResponse === 'object') imgSize = JSON.stringify(imgResponse).length;
				if (imgStatus === 200) sizes.push(imgSize);
				else throw new Error(`${imageMethod} returned status ${imgStatus}`);
			}
		}
		const [faviconStatus, faviconResponse] = await this.favicon();
		let faviconSize = -1;
		if (typeof faviconResponse === 'string') faviconSize = faviconResponse.length;
		else if (faviconResponse instanceof ArrayBuffer) faviconSize = faviconResponse.byteLength;
		else if (typeof faviconResponse === 'object') faviconSize = JSON.stringify(faviconResponse).length;
		if (faviconStatus === 200) sizes.push(faviconSize);
		else throw new Error(`favicon returned status ${faviconStatus}`);

		return sizes;
	}
}

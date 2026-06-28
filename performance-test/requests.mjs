//@ts-check
export class Requests {
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
		return [response.status, await response.json()];
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

	async api3() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/api/goalPost?goalId=247488&postDateTime=1782122907',
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

	async image1() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=1782122907&index=0',
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

	async image2() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=1782122907&index=1',
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

	async image3() {
		const response = await fetch(
			'http://192.168.0.23:30001/hinst-website/api/goalPost/image?goalId=247488&postDateTime=1782122907&index=2',
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
			'api2',
			'api3',
			'favicon',
			'image1',
			'image2',
			'image3'
		];
		const sizes = [];
		for (const name of methods) {
			const [status, response] = await this[name]();
			let responseInfo = -1;
			if (typeof response === 'string') responseInfo = response.length;
			else if (response instanceof ArrayBuffer) responseInfo = response.byteLength;
			else if (typeof response === 'object' && response !== null)
				responseInfo = Object.keys(response).length;
			if (status === 200) sizes.push(responseInfo);
			else throw new Error(`${name} returned status ${status}`);
		}
		return sizes;
	}
}

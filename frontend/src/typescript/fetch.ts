async function fetchSafe(url: string, requestInit: RequestInit = {}): Promise<Response> {
	const response = await fetch(url, requestInit);
	if (!response.ok) throw new Error(url + ' ' + response.statusText);
	return response;
}

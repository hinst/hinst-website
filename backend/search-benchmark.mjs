// @ts-check
const WORDS = ['pascal', 'always', 'week'];
const API_URL = 'http://192.168.0.23:30001/hinst-website/api';
const COUNT = 100;

/** @param {string} word */
async function benchmark(word) {
    const url = API_URL + '/goalPosts/search?query=' + encodeURIComponent(word);
    let totalTime = 0;
    /** @type {number|undefined} */
    let minTime = undefined;
    /** @type {number|undefined} */
    let maxTime = undefined;
    /** @type {number|undefined} */
    let resultCount = undefined;
    for (let i = 0; i < COUNT; i++) {
        const startTime = performance.now();
        const response = await fetch(url, { headers: { 'Accept-Language': 'en-GB' } });
        const data = (await response.json()) || [];
        const currentResultCount = data.length;
        if (resultCount == null)
            resultCount = currentResultCount;
        else
            console.assert(resultCount === currentResultCount,
                `Expected ${resultCount} results but got ${currentResultCount}`);
        const endTime = performance.now();
        const elapsedTime = endTime - startTime;
        totalTime += elapsedTime;
        if (minTime == null || elapsedTime < minTime)
            minTime = elapsedTime;
        if (maxTime == null || elapsedTime > maxTime)
            maxTime = elapsedTime;
    }
    minTime = minTime || 0;
    maxTime = maxTime || 0;
    resultCount = resultCount || 0;
    const avgTime = totalTime / COUNT;
    console.log(`Search for "${word}": min=${minTime.toFixed(0)}ms, avg=${avgTime.toFixed(0)}ms, max=${maxTime.toFixed(0)}ms, results=${resultCount}`);
}

async function main() {
    console.log(`Benchmarking: count=${COUNT}, url=${API_URL}, words=${WORDS.length}`);
    for (const word of WORDS)
        await benchmark(word);
}
main();
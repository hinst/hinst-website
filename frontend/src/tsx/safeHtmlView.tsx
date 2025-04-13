export default function SafeHtmlView(props: { htmlText: string }) {
	const parser = new DOMParser();
	const parsedHtml = parser.parseFromString(props.htmlText, 'text/html');
	const scriptTags = parsedHtml.getElementsByTagName('script');
	[...scriptTags].forEach((tag) => tag.remove());
	return <div dangerouslySetInnerHTML={{ __html: parsedHtml.body.innerHTML }} />;
}

export default function SafeHtmlView(props: {
	htmlText: string;
	updateDocument?: (document: Document) => void;
}) {
	const parser = new DOMParser();
	const parsedHtml = parser.parseFromString(props.htmlText, 'text/html');
	const scriptTags = parsedHtml.getElementsByTagName('script');
	[...scriptTags].forEach((tag) => tag.remove());
	if (props.updateDocument) props.updateDocument(parsedHtml);
	return <div dangerouslySetInnerHTML={{ __html: parsedHtml.body.innerHTML }} />;
}

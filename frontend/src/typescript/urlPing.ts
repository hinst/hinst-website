export class UrlPingRecord {
	constructor(
		public readonly url: string,
		public readonly googlePingedAt: number | null,
		public readonly googlePingedManuallyAt: number | null
	) {}
}

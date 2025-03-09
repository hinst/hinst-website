export class GoalHeader {
	constructor(
		public id: string,
		public title: string,
		public postCount: number,
		public updatedAt: string,
		public author: string,
	) {
	}
}

export class PostHeader {
	constructor(
		public id: string,
		public date: string,
		/** Goal id */
		public obj_id: string,
	) {
	}
}
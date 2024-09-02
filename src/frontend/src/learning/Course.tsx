export interface Course {
	id: string;
	title: string;
	description: string;
	topics: string[];
	difficulty: string;
	duration: string;
	language: string;
	isFree: boolean;
	url: string;
}
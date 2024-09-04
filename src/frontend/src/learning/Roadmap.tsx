// Define the Roadmap interface
import {Course} from "./Course.tsx";

export interface Roadmap {
	id: string;
	title: string;
	author: string;
	courses: Course[];
	courseIDs: string[];
	topics: string[];
	isCustom: boolean;
	createdBy: string;
	likes: number;
	difficulty: string;
	imageUrl: string;
	liked: boolean;
}
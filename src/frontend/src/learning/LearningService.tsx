import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService';
import {Roadmap} from "./Roadmap.tsx";

const GET_NAVBAR_DATA = gql`
    query GetNavbarData($name: String!) {
        getUserByName(name: $name) {
            dailyChallengeAvailable
            role
        }
    }
`;


const GET_COURSES = gql`
    query GetCourses($userId: String!, $pagination: PaginationInput!) {
        getCourses(userId: $userId, pagination: $pagination) {
            courses {
                id
                title
                url
                description
                difficulty
                topics
                isFree
                duration
                language
            }
            pagination {
                page
                perPage
                lastEvaluatedKey
            }
        }
    }
`;

const GET_ROADMAPS = gql`
    query GetRoadmaps {
        getRoadmaps {
            id
            title
            author
            topics
            isCustom
            likes
            difficulty
        }
    }
`;

const GET_ROADMAP_BY_ID = gql`
    query GetRoadmapById($id: ID!, $userId: String!, $from: String) {
        getRoadmapById(id: $id, userId: $userId, from: $from) {
            id
            title
            author
            courses {
                id
                title
                description
                difficulty
                duration
                url
            }
            topics
            isCustom
            likes
            difficulty
            description
        }
    }
`;

const UPSERT_ROADMAP = gql`
    mutation UpsertRoadmap($input: RoadmapInput!) {
        upsertRoadmap(input: $input) {
            id
        }
    }
`;

const UPSERT_COURSE = gql`
    mutation UpsertCourse($input: CourseInput!) {
        upsertCourse(input: $input) {
            id
        }
    }
`;

const USER_LIKED_ROADMAP = gql`
    mutation UserLikedRoadmap($userId: ID!, $roadmapId: ID!) {
        userLikedRoadmap(userId: $userId, roadmapId: $roadmapId) {
            success
        }
    }
`;

const GET_ROADMAPS_BY_USER = gql`
    query GetRoadmapsByUser($userId: String!) {
        getRoadmapsByUser(userId: $userId) {
            id
            title
            author
            topics
            isCustom
            likes
            difficulty
        }
    }
`;

const CUSTOM_ROADMAP_REQUESTED = gql`
    mutation CustomRoadmapRequested($prompt: String, $userId: String) {
        customRoadmapRequested(prompt: $prompt, userId: $userId) {
            title
            topics
            likes
            difficulty
            description
            courses {
                id
                title
                description
                difficulty
                duration
                url
                isFree
            }
        }
    }
`;

const GET_ROADMAP_FEED = gql`
    query GetRoadmapFeed($userId: String!) {
        getRoadmapFeed(userId: $userId) {
            id
            title
            author
            topics
            isCustom
            likes
            difficulty
            liked
            imageUrl
        }
    }
`;

class LearningService {
	private client: ApolloClient<any>;
	private authService: AuthService;

	constructor(client: ApolloClient<any>) {
		this.client = client;
		this.authService = new AuthService(client);
	}

	async getNavbarData(): Promise<{ dailyChallengeAvailable: boolean, role: number }> {
		const username = this.authService.getCognitoUsername();
		const {data} = await this.client.query({
			query: GET_NAVBAR_DATA,
			variables: {name: username},
		});
		return data.getUserByName;
	}

	async getCourses(userId: string, pagination: {
		page: number,
		perPage: number,
		lastEvaluatedKey?: any
	}): Promise<any> {
		let temp = {
			page: pagination.page,
			perPage: pagination.perPage,
			lastEvaluatedKey: pagination.lastEvaluatedKey
		};

		const {data} = await this.client.query({
			query: GET_COURSES,
			variables: {userId: userId, pagination: temp},
		});
		return data.getCourses;
	}

	async getRoadmaps(): Promise<any[]> {
		const {data} = await this.client.query({
			query: GET_ROADMAPS,
		});
		return data.getRoadmaps;
	}

	async getRoadmapById(id: string): Promise<Roadmap> {
		const authService = new AuthService(this.client);
		const userId = authService.getCognitoUsername()

		const {data} = await this.client.query({
			query: GET_ROADMAP_BY_ID,
			variables: {id, userId},
		});

		return data.getRoadmapById as Roadmap;
	}

	async upsertRoadmap(input: any): Promise<any> {
		const {data} = await this.client.mutate({
			mutation: UPSERT_ROADMAP,
			variables: {input},
		});
		return data.upsertRoadmap;
	}

	async upsertCourse(input: any): Promise<any> {
		const {data} = await this.client.mutate({
			mutation: UPSERT_COURSE,
			variables: {input},

		});
		return data.upsertCourse;
	}

	async userLikedRoadmap(userId: string, roadmapId: string): Promise<boolean> {
		const {data} = await this.client.mutate({
			mutation: USER_LIKED_ROADMAP,
			variables: {userId, roadmapId},
		});
		return data.userLikedRoadmap.success;
	}

	async getRoadmapsByUser(userId: string): Promise<any[]> {
		const {data} = await this.client.query({
			query: GET_ROADMAPS_BY_USER,
			variables: {userId},
		});
		return data.getRoadmapsByUser;
	}

	async getRoadmapFeed(userId: string): Promise<any> {
		const {data} = await this.client.query({
			query: GET_ROADMAP_FEED,
			variables: {userId},
			fetchPolicy: 'no-cache',
		});
		return data.getRoadmapFeed;
	}

	async customRoadmapRequested(prompt: string, userId: string): Promise<any> {
		const {data} = await this.client.mutate({
			mutation: CUSTOM_ROADMAP_REQUESTED,
			variables: {prompt, userId},
		});
		return data.customRoadmapRequested;
	}
}

export default LearningService;
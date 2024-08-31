import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService';

const GET_NAVBAR_DATA = gql`
    query GetNavbarData($name: String!) {
        getUserByName(name: $name) {
            dailyChallengeAvailable
            role
        }
    }
`;

const GET_COURSES = gql`
    query GetCourses {
        getCourses {
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
    query GetRoadmapById($id: ID!) {
        getRoadmapById(id: $id) {
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
    mutation CustomRoadmapRequested($prompt: String) {
        customRoadmapRequested(prompt: $prompt) {
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

	async getCourses(): Promise<any[]> {
		const {data} = await this.client.query({
			query: GET_COURSES,
		});
		return data.getCourses;
	}

	async getRoadmaps(): Promise<any[]> {
		const {data} = await this.client.query({
			query: GET_ROADMAPS,
		});
		return data.getRoadmaps;
	}

	async getRoadmapById(id: string): Promise<any> {
		const {data} = await this.client.query({
			query: GET_ROADMAP_BY_ID,
			variables: {id},
		});
		return data.getRoadmapById;
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
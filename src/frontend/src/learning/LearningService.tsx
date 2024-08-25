import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService';

const GET_DAILY_CHALLENGE_AVAILABILITY = gql`
    query GetDailyChallengeAvailability($name: String!) {
        getUserByName(name: $name) {
            dailyChallengeAvailable
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



class LearningService {
	private client: ApolloClient<any>;
	private authService: AuthService;

	constructor(client: ApolloClient<any>) {
		this.client = client;
		this.authService = new AuthService(client);
	}

	async getDailyChallengeAvailability(): Promise<boolean> {
		const username = this.authService.getCognitoUsername();
		const {data} = await this.client.query({
			query: GET_DAILY_CHALLENGE_AVAILABILITY,
			variables: {name: username},
		});
		return data.getUserByName.dailyChallengeAvailable;
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
}

export default LearningService;
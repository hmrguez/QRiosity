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
            createdBy
            likes
            difficulty
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
}

export default LearningService;
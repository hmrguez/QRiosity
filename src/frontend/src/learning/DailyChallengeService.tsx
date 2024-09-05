import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService.tsx';

const GET_DAILY_CHALLENGE = gql`
    query GetDailyChallenge($userId: String!) {
        dailyChallenge(userId: $userId) {
            question
            categories
            type
        }
    }
`;

const SUBMIT_DAILY_CHALLENGE = gql`
    mutation SubmitDailyChallenge($username: String!, $question: String!, $answer: String!) {
        dailyChallenge(username: $username, question: $question, answer: $answer) {
            userId
            question
            answer
            rating
            insight
            left
        }
    }
`;

class DailyChallengeService {
	private client: ApolloClient<any>;
	private authService: AuthService;

	constructor(client: ApolloClient<any>) {
		this.client = client;
		this.authService = new AuthService(client);
	}

	async getDailyChallenge(): Promise<any> {
		const username = this.authService.getCognitoUsername();
		const {data} = await this.client.query({
			query: GET_DAILY_CHALLENGE,
			variables: {userId: username},
			fetchPolicy: 'no-cache',
		});
		return data.dailyChallenge;
	}

	async submitDailyChallenge(question: string, answer: string): Promise<any> {
		const username = this.authService.getCognitoUsername();
		const {data} = await this.client.mutate({
			mutation: SUBMIT_DAILY_CHALLENGE,
			variables: {username, question, answer},
		});
		return data.dailyChallenge;
	}
}

export default DailyChallengeService;
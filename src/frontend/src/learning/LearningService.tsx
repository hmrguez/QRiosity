import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService';

const GET_DAILY_CHALLENGE_AVAILABILITY = gql`
    query GetDailyChallengeAvailability($name: String!) {
        getUserByName(name: $name) {
            dailyChallengeAvailable
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

		const username = this.authService.getUsername();

		const {data} = await this.client.query({
			query: GET_DAILY_CHALLENGE_AVAILABILITY,
			variables: {name: username},
		});

		return data.getUserByName.dailyChallengeAvailable;
	}


}

export default LearningService;
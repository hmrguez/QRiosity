// src/frontend/src/services/ProblemService.ts
import {ApolloClient, gql} from '@apollo/client';

class ProblemService {
	private client: ApolloClient<any>;

	constructor(client: ApolloClient<any>) {
		this.client = client;
	}

	async getDailyChallenge(userId: string): Promise<any> {
        const query = gql`
            query {
                dailyChallenge(userId: "${userId}") {
                    question
                    categories
                    type
                }
            }
		`;

		const response = await this.client.query({query});
		return response.data.dailyChallenge;
    }
}

export default ProblemService;
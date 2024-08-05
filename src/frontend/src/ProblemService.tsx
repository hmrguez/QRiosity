// src/frontend/src/services/ProblemService.ts
import { gql, ApolloClient } from '@apollo/client';

class ProblemService {
  private client: ApolloClient<any>;

  constructor(client: ApolloClient<any>) {
    this.client = client;
  }

  async getDailyChallenge(category: string): Promise<any> {
    const query = gql`
        query {
            dailyChallenge(category: "${category}") {
                id
                question
                categories
                type
            }
        }
	`;

    const response = await this.client.query({ query });
    return response.data.dailyChallenge;
  }

  // async rateDailyChallenge(challengeId: string, rating: number): Promise<any> {
  //   const mutation = gql`
  //     mutation($challengeId: ID!, $rating: Int!) {
  //       dailyChallenge(userId: $challengeId, rating: $rating) {
  //         success
  //         message
  //       }
  //     }
  //   `;
  //
  //   const variables = { challengeId, rating };
  //   const response = await this.client.mutate({ mutation, variables });
  //   return response.data.rateDailyChallenge;
  // }
}

export default ProblemService;
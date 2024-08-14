import { ApolloClient, gql } from '@apollo/client';

class TopicService {
    private client: ApolloClient<any>;

    constructor(client: ApolloClient<any>) {
        this.client = client;
    }

    async getAllTopics(): Promise<string[]> {
        const query = gql`
            query {
                getAllTopics {
                    name
                }
            }
        `;

        const response = await this.client.query({ query });
		return response.data.getAllTopics.map((topic: { name: string }) => topic.name)
    }
}

export default TopicService;
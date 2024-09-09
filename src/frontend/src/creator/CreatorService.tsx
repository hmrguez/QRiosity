import {ApolloClient, gql} from '@apollo/client';
import AuthService from '../auth/AuthService';

const GET_USER_BY_NAME = gql`
    query GetUserByName($name: String!) {
        getUserByName(name: $name) {
            genUsagesRemaining
            creationsRemaining
            role
        }
    }
`;

class CreatorService {
	private client: ApolloClient<any>;
	private authService: AuthService;

	constructor(client: ApolloClient<any>) {
		this.client = client;
		this.authService = new AuthService(client);
	}

	async getBuilderStats(): Promise<{ usages: number, role: number }> {
		const username = this.authService.getCognitoUsername();
		const {data} = await this.client.query({
			query: GET_USER_BY_NAME,
			variables: {name: username},
		});
		return {
			usages: data.getUserByName.genUsagesRemaining,
			role: data.getUserByName.role,
		};
	}
}

export default CreatorService;
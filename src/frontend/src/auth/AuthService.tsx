import {ApolloClient, gql} from '@apollo/client';
import {jwtDecode} from "jwt-decode";

const LOGIN_QUERY = gql`
    query Login($username: String!, $password: String!) {
        login(username: $username, password: $password) {
            token
            user {
                id
                name
            }
        }
    }
`;

const REGISTER_MUTATION = gql`
    mutation Register($username: String!, $password: String!, $email: String!, $topics: [String!]!) {
        register(username: $username, password: $password, email: $email, topics: $topics) {
            token
            user {
                id
                name
            }
        }
    }
`;

class AuthService {
	client;


	constructor(client: ApolloClient<any>) {
		this.client = client;
	}

	async login(username: string, password: string) {
		const {data} = await this.client.query({
			query: LOGIN_QUERY,
			variables: {username, password},
		});
		localStorage.setItem('token', data.login.token);
		return data.login;
	}

	async register(username: string, password: string, email: string, topics: string[]) {
		const {data} = await this.client.mutate({
			mutation: REGISTER_MUTATION,
			variables: {username, password, email, topics},
		});
		localStorage.setItem('token', data.register.token);
		return data.register;
	}

	isLoggedIn() {
		return !!localStorage.getItem('token');
	}

	getUsername() {
		const token = localStorage.getItem('token');
		if (!token) {
			return null;
		}

		const decoded: { username: string, userId: string, exp: Date } = jwtDecode(token);

		return decoded.username;
	}

	logout() {
		localStorage.removeItem('token');
	}
}

export default AuthService;
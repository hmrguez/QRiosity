import {ApolloClient, gql} from '@apollo/client';
import {jwtDecode} from "jwt-decode";

const LOGIN_QUERY = gql`
    query Login($username: String!, $password: String!) {
        login(username: $username, password: $password) {
            token
        }
    }
`;

const REGISTER_MUTATION = gql`
    mutation Register($username: String!, $password: String!, $email: String!, $topics: [String!]!) {
        register(username: $username, password: $password, email: $email, topics: $topics) {
            token
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


	// TODO: UserId is not being sent correctly in the token
	getUserId() {
		const token = localStorage.getItem('token');
		if (!token) {
			return null;
		}

		const decoded: { username: string, userId: string, exp: Date } = jwtDecode(token);

		return decoded.userId;
	}

	getCognitoUsername() {
		const token = localStorage.getItem('token');
		if (!token) {
			return null;
		}

		const decoded: { sub: String } = jwtDecode(token);

		return decoded.sub;
	}

	getUsername() {
		const token = localStorage.getItem('token');
		if (!token) {
			return null;
		}

		const decoded: { name: String } = jwtDecode(token);

		return decoded.name;
	}

	logout() {
		localStorage.removeItem('token');
	}
}

export default AuthService;
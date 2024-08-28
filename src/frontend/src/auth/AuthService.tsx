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
        register(username: $username, password: $password, email: $email, topics: $topics)
    }
`;


const CONFIRM_EMAIL_MUTATION = gql`
    mutation ConfirmEmail($email: String!, $token: String!) {
        confirmEmail(email: $email, token: $token){
            success
        }
    }
`;

const GET_USER_BY_NAME_QUERY = gql`
    query GetUserByName($name: String!) {
        getUserByName(name: $name) {
            username
            topics
            role
        }
    }
`;

const UPDATE_USER_MUTATION = gql`
    mutation UpdateUser($input: UserEditInput!) {
        updateUser(input: $input) {
            success
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
		return data.register;
	}

	async confirmEmail(email: string, token: string): Promise<boolean> {
		const response = await this.client.mutate({
			mutation: CONFIRM_EMAIL_MUTATION,
			variables: {email, token}
		});
		return response.data.confirmEmail;
	}

	async resendConfirmationEmail(email: string): Promise<boolean> {
        const query = gql`
            query ResendConfirmationEmail($email: String!) {
                resendConfirmationEmail(email: $email){
                    success
                }
            }
		`;
		const response = await this.client.query({
			query,
			variables: {email}
		});
		return response.data.resendConfirmationEmail;
    }

	async getProfile(name: string) {
		const {data} = await this.client.query({
			query: GET_USER_BY_NAME_QUERY,
			variables: {name},
		});
		return data.getUserByName;
	}

	async updateUser(input: {
		name: string;
		role: number;
		email: string;
		topics: string[];
		dailyChallengeAvailable: boolean;
	}) {
		const {data} = await this.client.mutate({
			mutation: UPDATE_USER_MUTATION,
			variables: {input},
		});
		return data.updateUser;
	}


	isLoggedIn() {
		return !!localStorage.getItem('token');
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
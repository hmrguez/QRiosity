// src/utils/apolloLink.ts
import {ApolloLink} from '@apollo/client';

export const authLink = new ApolloLink((operation, forward) => {

	const token = localStorage.getItem('token');
	operation.setContext(({headers = {}}) => ({
		headers: {
			...headers,
			Authorization: token ? `Bearer ${token}` : '',
		}
	}));
	return forward(operation);
});
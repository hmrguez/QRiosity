// src/frontend/src/apolloClient.ts
import {ApolloClient, HttpLink, InMemoryCache} from '@apollo/client';

const client = new ApolloClient({
	link: new HttpLink({
		uri: 'https://4sxml7ptwfb7xpbyjbu24jbpnu.appsync-api.us-east-2.amazonaws.com/graphql',  // Replace with your AppSync API URL
		headers: {
			'x-api-key': "da2-h4ddio3hgnbtvoz4p5vswbp3fe",  // Replace with your AppSync API key if using API key authentication
		},
	}),
	cache: new InMemoryCache(),
});


// const client = new ApolloClient({
// 	uri: 'http://localhost:9000/query',
// 	cache: new InMemoryCache(),
// });

export default client;


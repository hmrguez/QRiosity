import {useApolloClient} from '@apollo/client';

const SomeComponent = () => {
	const client = useApolloClient();

	const clearCache = async () => {
		try {
			await client.clearStore(); // Clears the cache and refetches active queries
			// or use client.resetStore(); // Clears the cache and refetches active queries
			console.log('Cache cleared');
		} catch (error) {
			console.error('Error clearing cache', error);
		}
	};

	return (
		<button onClick={clearCache}>Clear Cache</button>
	);
};

export default SomeComponent;
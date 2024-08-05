// src/frontend/src/components/DailyChallenge.tsx
import React from 'react';
import {gql, useQuery} from '@apollo/client';

const GET_DAILY_CHALLENGE = gql`
    query {
        dailyChallenge(category: "javascript") {
            id
            question
            categories
            type
        }
    }
`;

const DailyChallenge: React.FC = () => {
	const {loading, error, data} = useQuery(GET_DAILY_CHALLENGE);

	if (loading) return <div>Loading...</div>;
	if (error) return <div>Error: {error.message}</div>;

	const {dailyChallenge} = data;

	return (
		<div>
			<h1>{dailyChallenge.question}</h1>
			<p>Categories: {dailyChallenge.categories.join(', ')}</p>
			<p>Type: {dailyChallenge.type}</p>
		</div>
	);
};

export default DailyChallenge;
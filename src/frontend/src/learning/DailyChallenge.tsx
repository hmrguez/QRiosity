// src/frontend/src/components/DailyChallenge.tsx
import React, {useState} from 'react';
import {gql, useMutation, useQuery} from '@apollo/client';
import {Button} from 'primereact/button';
import {InputText} from 'primereact/inputtext';

interface DailyChallengeResult {
	userId: string;
	question: string;
	answer: string;
	rating: number;
	insight: string;
}

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

const SUBMIT_DAILY_CHALLENGE = gql`
    mutation SubmitDailyChallenge($userId: ID!, $question: String!, $answer: String!) {
        dailyChallenge(userId: $userId, question: $question, answer: $answer) {
            userId
            question
            answer
            rating
            insight
        }
    }
`;

interface DailyChallengeProps {
	answer: string;
	setAnswer: (answer: string) => void;
}

const DailyChallenge: React.FC<DailyChallengeProps> = ({answer, setAnswer}) => {
	const {loading, error, data} = useQuery(GET_DAILY_CHALLENGE);
	const [submitDailyChallenge] = useMutation(SUBMIT_DAILY_CHALLENGE);
	const [result, setResult] = useState<DailyChallengeResult | null>(null);

	if (loading) return <div>Loading...</div>;
	if (error) return <div>Error: {error.message}</div>;

	const {dailyChallenge} = data;

	const handleSubmit = async () => {
		try {
			const {data} = await submitDailyChallenge({
				variables: {
					userId: '1', // Placeholder userId
					question: dailyChallenge.question,
					answer: answer
				}
			});
			setResult(data.dailyChallenge);
		} catch (error) {
			console.error('Error submitting challenge:', error);
		}
	};

	return (
		<div>
			<h1>{dailyChallenge.question}</h1>
			<p>Categories: {dailyChallenge.categories.join(', ')}</p>
			<p>Type: {dailyChallenge.type}</p>
			<div className="p-field">
				<label htmlFor="answer">Your Answer</label>
				<InputText id="answer" value={answer} onChange={(e) => setAnswer(e.target.value)}/>
			</div>
			<Button label="Submit" onClick={handleSubmit}/>
			{result && <h2>Rating: {result.rating}</h2>}
			{result && <h2>Insight: {result.insight}</h2>}
		</div>
	);
};

export default DailyChallenge;
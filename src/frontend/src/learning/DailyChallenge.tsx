import React, {useEffect, useState} from 'react';
import {gql, useApolloClient, useMutation, useQuery} from '@apollo/client';
import {Button} from 'primereact/button';
import {InputTextarea} from "primereact/inputtextarea";
import './DailyChallenge.css';
import AuthService from "../auth/AuthService.tsx";

interface DailyChallengeResult {
	userId: string;
	question: string;
	answer: string;
	rating: number;
	insight: string;
	left: number;
}

const SUBMIT_DAILY_CHALLENGE = gql`
    mutation SubmitDailyChallenge($username: String!, $question: String!, $answer: String!) {
        dailyChallenge(username: $username, question: $question, answer: $answer) {
            userId
            question
            answer
            rating
            insight
            left
        }
    }
`;

interface DailyChallengeProps {
	answer: string;
	setAnswer: (answer: string) => void;
	onHide: () => void;
	onCorrectSubmit: () => void;
}

const DailyChallenge: React.FC<DailyChallengeProps> = ({answer, setAnswer, onHide, onCorrectSubmit}) => {
	const client = useApolloClient();
	const authService = new AuthService(client);
	const username = authService.getCognitoUsername();

    const GET_DAILY_CHALLENGE = gql`
        query {
            dailyChallenge(userId: "${username}") {
                question
                categories
                type
            }
        }
	`;

	const {loading, error, data} = useQuery(GET_DAILY_CHALLENGE);
	const [submitDailyChallenge] = useMutation(SUBMIT_DAILY_CHALLENGE);
	const [result, setResult] = useState<DailyChallengeResult | null>(null);
	const [left, setLeft] = useState<number>(0);

	const [buttonClass, setButtonClass] = useState('p-button-contrast');
	const [resultLoading, setResultLoading] = useState(false);
	const [buttonIcon, setButtonIcon] = useState('pi pi-check');
	const [buttonLabel, setButtonLabel] = useState('Submit');

	useEffect(() => {
		if (result) {
			setResultLoading(false);
			setLeft(result.left);
			if (result.rating >= 6) {
				setButtonClass('p-button-success');
				setButtonIcon('pi pi-thumbs-up');
				setButtonLabel('Success');
			} else {
				setButtonClass('p-button-danger');
				setButtonIcon('pi pi-thumbs-down');
				setButtonLabel('Failed');
			}
		}
	}, [result]);

	if (loading) return <div>Loading...</div>;
	if (error) return <div>Error: {error.message}</div>;

	const {dailyChallenge} = data;

	const handleExit = () => {
		setAnswer('');
		onHide();
	};

	const handleSubmit = async () => {
		setResultLoading(true);

		try {
			const {data} = await submitDailyChallenge({
				variables: {
					username: authService.getCognitoUsername(),
					question: dailyChallenge.question,
					answer: answer
				}
			});
			onCorrectSubmit();
			setResult(data.dailyChallenge);
		} catch (error) {
			console.error('Error submitting challenge:', error);
		}
	};

	const handleNext = () => {
		setResult(null);
		setAnswer('');
		setButtonClass('p-button-contrast');
		setButtonIcon('pi pi-check');
		setButtonLabel('Submit');
	};

	return (
		<div className="challenge-content">
			<h2 className="question-title">{dailyChallenge.question}</h2>
			<InputTextarea
				value={answer}
				onChange={(e) => setAnswer(e.target.value)}
				rows={5}
				cols={30}
				className="textarea"
				autoResize
				placeholder="Type your answer here..."
			/>
			<div className="button-group">
				<Button
					label="Exit"
					icon="pi pi-times"
					className="p-button-secondary daily-exit-button"
					onClick={handleExit}
				/>
				<Button
					disabled={!!result}
					label={buttonLabel}
					icon={buttonIcon}
					className={buttonClass + " daily-submit-button"}
					onClick={handleSubmit}
					loading={resultLoading}
				/>
				{result && left > 0 && (
					<Button
						label="Next"
						icon="pi pi-arrow-right"
						className="p-button-info daily-next-button"
						onClick={handleNext}
					/>
				)}
			</div>
			{result && (
				<div className="insight-content">
					<p className="insight-text">{result.insight}</p>
				</div>
			)}
		</div>
	);
};

export default DailyChallenge;
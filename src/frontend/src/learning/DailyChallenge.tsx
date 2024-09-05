import React, {useEffect, useState} from 'react';
import {useApolloClient} from '@apollo/client';
import {Button} from 'primereact/button';
import {InputTextarea} from "primereact/inputtextarea";
import './DailyChallenge.css';
import DailyChallengeService from './DailyChallengeService.tsx';

interface DailyChallengeResult {
	userId: string;
	question: string;
	answer: string;
	rating: number;
	insight: string;
	left: number;
}

interface DailyChallengeProps {
	answer: string;
	setAnswer: (answer: string) => void;
	onHide: () => void;
	onCorrectSubmit: () => void;
}

const DailyChallenge: React.FC<DailyChallengeProps> = ({answer, setAnswer, onHide, onCorrectSubmit}) => {
	const client = useApolloClient();
	const dailyChallengeService = new DailyChallengeService(client);

	const [dailyChallenge, setDailyChallenge] = useState<any>(null);
	const [result, setResult] = useState<DailyChallengeResult | null>(null);
	const [left, setLeft] = useState<number>(0);
	const [isLoading, setIsLoading] = useState<boolean>(true);

	const [buttonClass, setButtonClass] = useState('p-button-contrast');
	const [resultLoading, setResultLoading] = useState(false);
	const [buttonIcon, setButtonIcon] = useState('pi pi-check');
	const [buttonLabel, setButtonLabel] = useState('Submit');

	useEffect(() => {
		const fetchDailyChallenge = async () => {
			setIsLoading(true);
			try {
				const challenge = await dailyChallengeService.getDailyChallenge();
				setDailyChallenge(challenge);
			} catch (error) {
				console.error('Error fetching daily challenge:', error);
			} finally {
				setIsLoading(false);
			}
		};

		fetchDailyChallenge();
	}, []);

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

	if (isLoading) return <div>Loading...</div>;

	const handleExit = () => {
		setAnswer('');
		onHide();
	};

	const handleSubmit = async () => {
		setResultLoading(true);

		try {
			const challengeResult = await dailyChallengeService.submitDailyChallenge(dailyChallenge.question, answer);
			onCorrectSubmit();
			setResult(challengeResult);
		} catch (error) {
			console.error('Error submitting challenge:', error);
		}
	};

	const handleNext = async () => {
		setIsLoading(true);
		try {
			const newChallenge = await dailyChallengeService.getDailyChallenge();
			setDailyChallenge(newChallenge);
			setResult(null);
			setAnswer('');
			setButtonClass('p-button-contrast');
			setButtonIcon('pi pi-check');
			setButtonLabel('Submit');
		} catch (error) {
			console.error('Error fetching new challenge:', error);
		} finally {
			setIsLoading(false);
		}
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
						className="p-button-contrast daily-submit-button"
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
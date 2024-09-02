import React, {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useApolloClient} from '@apollo/client';
import AuthService from './AuthService';
import './Login.css';
import {useAuth} from './AuthContext.tsx';
import {MultiSelect} from 'primereact/multiselect';
import {Steps} from 'primereact/steps';
import TopicService from "../learning/TopicService.tsx";

const RegisterForm = () => {
	const [username, setUsername] = useState('');
	const [email, setEmail] = useState('');
	const [password, setPassword] = useState('');
	const [confirmPassword, setConfirmPassword] = useState('');
	const [topics, setTopics] = useState<string[]>([]);
	const [error, setError] = useState('');
	const [activeIndex, setActiveIndex] = useState(0);
	const [allTopics, setAllTopics] = useState([])

	const client = useApolloClient();
	const authService = new AuthService(client);
	const topicService = new TopicService(client)
	const navigate = useNavigate();
	const {login} = useAuth();

	// Useeffect to get allTopics
	useEffect(() => {
		const fetchTopics = async () => {
			const temp = await topicService.getAllTopics();
			// @ts-ignore
			setAllTopics(temp);
		};
		fetchTopics();
	}, []);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (password !== confirmPassword) {
			setError('Passwords do not match.');
			return;
		}
		try {
			await authService.register(username, password, email, topics);
			login();
			navigate('/confirm-email' + "?email=" + email);
		} catch (err: any) {
			console.error(err);
			setError(err.message);
		}
	};

	const items = [
		{label: 'User Info'},
		{label: 'Topics'},
		{label: 'Password'}
	];

	return (
		<div className="login-container">
			<div className="logo">Register</div>
			<Steps className="mb-4" model={items} activeIndex={activeIndex} onSelect={(e) => setActiveIndex(e.index)}/>
			<form onSubmit={handleSubmit}>
				{activeIndex === 0 && (
					<>
						<label htmlFor="username">Username</label>
						<input
							type="text"
							id="username"
							name="username"
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							required
						/>

						<label htmlFor="email">Email</label>
						<input
							type="email"
							id="email"
							name="email"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							required
						/>
					</>
				)}

				{activeIndex === 1 && (
					<>
						<label htmlFor="topics">Topics</label>
						<MultiSelect id="topics" value={topics}
									 options={allTopics.map(topic => ({label: topic, value: topic}))}
									 onChange={(e) => setTopics(e.value)} className="w-full mb-3"/>
					</>
				)}

				{activeIndex === 2 && (
					<>
						<label htmlFor="password">Password</label>
						<input
							type="password"
							id="password"
							name="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							required
						/>

						<label htmlFor="confirmPassword">Confirm Password</label>
						<input
							type="password"
							id="confirmPassword"
							name="confirmPassword"
							value={confirmPassword}
							onChange={(e) => setConfirmPassword(e.target.value)}
							required
						/>
					</>
				)}

				{error && <div className="error-message">{error}</div>}

				<div className="button-row mt-4">
					{activeIndex > 0 && (
						<button type="button" className="back-button" onClick={() => setActiveIndex(activeIndex - 1)}>
							Back
						</button>
					)}
					{activeIndex < 2 && (
						<button type="button" className="next-button" onClick={() => setActiveIndex(activeIndex + 1)}>
							Next
						</button>
					)}
					{activeIndex === 2 && (
						<button type="submit" className="register-button">
							Register
						</button>
					)}
				</div>
			</form>
		</div>
	);
};

export default RegisterForm;
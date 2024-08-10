import {useEffect, useRef, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {gql, useApolloClient} from '@apollo/client';
import AuthService from './AuthService';
import {InputText} from 'primereact/inputtext';
import {MultiSelect} from 'primereact/multiselect';
import {Password} from 'primereact/password';
import {Button} from 'primereact/button';
import {Steps} from 'primereact/steps';
import {Toast} from 'primereact/toast';
import {useAuth} from "./AuthContext.tsx";

const GET_ALL_TOPICS = gql`
    query GetAllTopics {
        getAllTopics {
            name
        }
    }
`;

const ADD_TOPICS = gql`
    mutation AddTopics($names: [String!]!) {
        addTopics(names: $names) {
            name
        }
    }
`;

const Register = () => {
	const [activeIndex, setActiveIndex] = useState(0);
	const [username, setUsername] = useState('');
	const [email, setEmail] = useState('');
	const [topics, setTopics] = useState([]);
	const [allTopics, setAllTopics] = useState([]);
	const [newTopics, setNewTopics] = useState([]);
	const [password, setPassword] = useState('');
	const [confirmPassword, setConfirmPassword] = useState('');
	const client = useApolloClient();
	const authService = new AuthService(client);
	const navigate = useNavigate();
	const toast = useRef<Toast>(null);
	const { login } = useAuth();


	useEffect(() => {
		const fetchTopics = async () => {
			const {data} = await client.query({query: GET_ALL_TOPICS});
			setAllTopics(data.getAllTopics.map((topic: { name: string }) => topic.name));
		};
		fetchTopics();
	}, [client]);

	const handleNext = () => {
		if (activeIndex < steps.length - 1) {
			setActiveIndex(activeIndex + 1);
		}
	};

	const handlePrev = () => {
		if (activeIndex > 0) {
			setActiveIndex(activeIndex - 1);
		}
	};

	const handleTopicChange = (e: any) => {
		const selectedTopics = e.value;
		const addedTopics = selectedTopics.filter((topic: never) => !allTopics.includes(topic));
		setTopics(selectedTopics);
		setNewTopics(addedTopics);
	};

	const handleSubmit = async () => {
		console.log('Register');
		if (password !== confirmPassword) {
			toast.current?.show({severity: 'error', summary: 'Error', detail: 'Passwords do not match'});
			return;
		}
		if (newTopics.length > 0) {
			await client.mutate({mutation: ADD_TOPICS, variables: {names: newTopics}});
		}

		await authService.register(username, password, email, topics)
			.then(() => {
				login()
				navigate('/')
			})
			.catch((error) => {
				toast.current?.show({severity: 'error', summary: 'Error', detail: error.message});
			})
	};

	const steps = [
		{label: 'Account'},
		{label: 'Topics'},
		{label: 'Password'}
	];

	return (
		<div className="flex align-items-center justify-content-center">
			<div className="surface-card p-4 shadow-2 border-round w-full lg:w-6">
				<Toast ref={toast}/>
				<Steps model={steps} activeIndex={activeIndex}/>
				<div>
					{activeIndex === 0 && (
						<div>
							<label htmlFor="username" className="block text-900 font-medium mb-2">Username</label>
							<InputText id="username" value={username} onChange={(e) => setUsername(e.target.value)}
									   className="w-full mb-3"/>
							<label htmlFor="email" className="block text-900 font-medium mb-2">Email</label>
							<InputText id="email" value={email} onChange={(e) => setEmail(e.target.value)}
									   className="w-full mb-3"/>
						</div>
					)}
					{activeIndex === 1 && (
						<div>
							<label htmlFor="topics" className="block text-900 font-medium mb-2">Topics</label>
							<MultiSelect id="topics" value={topics}
										 options={allTopics.map(topic => ({label: topic, value: topic}))}
										 onChange={handleTopicChange} className="w-full mb-3"/>
						</div>
					)}
					{activeIndex === 2 && (
						<div>
							<label htmlFor="password" className="block text-900 font-medium mb-2">Password</label>
							<Password id="password" value={password} onChange={(e) => setPassword(e.target.value)}
									  className="w-full mb-3"/>
							<label htmlFor="confirmPassword" className="block text-900 font-medium mb-2">Confirm
								Password</label>
							<Password id="confirmPassword" value={confirmPassword}
									  onChange={(e) => setConfirmPassword(e.target.value)} className="w-full mb-3"/>
						</div>
					)}
					<div className="flex justify-content-between mt-4">
						<Button label="Previous" icon="pi pi-chevron-left" onClick={handlePrev}
								disabled={activeIndex === 0}/>
						{activeIndex < steps.length - 1 && (
							<Button label="Next" icon="pi pi-chevron-right" onClick={handleNext}/>
						)}
						{activeIndex === steps.length - 1 && (
							<Button label="Register" icon="pi pi-check" onClick={handleSubmit}/>
						)}
					</div>
				</div>
			</div>
		</div>
	);
};

export default Register;
import {useState} from 'react';
import {Link, useNavigate} from 'react-router-dom';
import {useApolloClient} from '@apollo/client';
import AuthService from './AuthService';
import {InputText} from "primereact/inputtext";
import {Button} from "primereact/button";
import {useAuth} from "./AuthContext.tsx";

const Login = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const client = useApolloClient();
	const authService = new AuthService(client);
	const navigate = useNavigate();
	const { login } = useAuth();

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		await authService.login(username, password);
		login();
		navigate('/');
	};

	return (
		<div className="flex align-items-center justify-content-center">
			<div className="surface-card p-4 shadow-2 border-round w-full lg:w-6">
				<div className="text-center mb-5">
					<img src="/demo/images/blocks/logos/hyper.svg" alt="hyper" height={50} className="mb-3"/>
					<div className="text-900 text-3xl font-medium mb-3">Welcome Back</div>
					<span className="text-600 font-medium line-height-3">Don't have an account?</span>
					<Link to="/register" className="font-medium no-underline ml-2 text-blue-500 cursor-pointer">Register
						now!</Link>
				</div>

				<div>
					<label htmlFor="email" className="block text-900 font-medium mb-2">Email</label>
					<InputText value={username} onChange={e => setUsername(e.target.value)} id="email" type="text"
							   placeholder="Email address" className="w-full mb-3"/>

					<label htmlFor="password" className="block text-900 font-medium mb-2">Password</label>
					<InputText value={password} onChange={e => setPassword(e.target.value)} id="password"
							   type="password" placeholder="Password" className="w-full mb-3"/>

					<Button label="Sign In" icon="pi pi-user" className="w-full" onClick={handleSubmit}/>
				</div>
			</div>
		</div>
	);
};

export default Login;
import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useApolloClient } from '@apollo/client';
import AuthService from './AuthService';
import { useAuth } from './AuthContext.tsx';
import './Login.css'

const LoginForm = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
	const [error, setError] = useState('');

	const client = useApolloClient();
    const authService = new AuthService(client);
    const navigate = useNavigate();
    const { login } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
		try {
			await authService.login(username, password);
			login();
			navigate('/');
		} catch (err) {
			console.error(err);
			setError('Login failed. Please check your username and password.');
		}
    };

    return (
        <div className="login-container">
            <div className="logo">Login</div>
            <form onSubmit={handleSubmit}>
                <label htmlFor="username">Username</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />

                <label htmlFor="password">Password</label>
                <input
                    type="password"
                    id="password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />

				{error && <div className="error-message">{error}</div>}


				<button type="submit" className="login-button">Log In</button>
            </form>
            <div className="divider">
                <span>or</span>
            </div>
            <div className="google-wrapper">
                <button className="google-login">
                    <img src="../assets/icons8-google.svg" alt="Google logo"/>
                    Sign in with Google
                </button>
            </div>
        </div>
    );
};

export default LoginForm;
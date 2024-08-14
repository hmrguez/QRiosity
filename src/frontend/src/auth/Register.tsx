import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useApolloClient } from '@apollo/client';
import AuthService from './AuthService';
import image from '../assets/icons8-google.svg';
import './Login.css';
import {useAuth} from "./AuthContext.tsx";

const RegisterForm = () => {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');

    const client = useApolloClient();
    const authService = new AuthService(client);
    const navigate = useNavigate();

    const { login } = useAuth();


    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (password !== confirmPassword) {
            setError('Passwords do not match.');
            return;
        }
        try {
            await authService.register(username, password, email, []);
            login();
            navigate('/confirm-email' + "?email=" + email);
        } catch (err) {
            console.error(err);
            setError(err.message);
        }
    };

    return (
        <div className="login-container">
            <div className="logo">Register</div>
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

                <label htmlFor="email">Email</label>
                <input
                    type="email"
                    id="email"
                    name="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
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

                <label htmlFor="confirmPassword">Confirm Password</label>
                <input
                    type="password"
                    id="confirmPassword"
                    name="confirmPassword"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                />

                {error && <div className="error-message">{error}</div>}

                <button type="submit" className="login-button">Register</button>
            </form>
            <div className="divider">
                <span>or</span>
            </div>
            <div className="google-wrapper">
                <button className="google-login">
                    <img src={image} alt="Google logo" />
                    Sign up with Google
                </button>
            </div>
        </div>
    );
};

export default RegisterForm;
import React, { useRef, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import AuthService from './AuthService';
import './ConfirmEmail.css';
import { useApolloClient } from '@apollo/client';
import { Toast } from 'primereact/toast';

const ConfirmEmail = () => {
    const [code, setCode] = useState('');
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const email = searchParams.get('email');

    const client = useApolloClient();
    const authService = new AuthService(client);
    const toast = useRef<Toast>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (email) {
            const success = await authService.confirmEmail(email, code);
            if (success) {
                navigate('/login');
            } else {
                toast.current?.show({ severity: 'error', summary: 'Error', detail: 'Email confirmation failed', life: 3000 });
            }
        }
    };

    const handleResend = async () => {
        if (email) {
            const success = await authService.resendConfirmationEmail(email);
            if (success) {
                toast.current?.show({ severity: 'success', summary: 'Success', detail: 'Confirmation email resent', life: 3000 });
            } else {
                toast.current?.show({ severity: 'error', summary: 'Error', detail: 'Failed to resend confirmation email', life: 3000 });
            }
        }
    };

    return (
        <div className="confirm-email-container">
            <Toast ref={toast} />
            <div className="logo">Confirm Email</div>
            <p className="mt-2 mb-5 text-center">A confirmation email was sent to {email}. Please enter the code below to confirm your email address.</p>
            <form onSubmit={handleSubmit}>
                <label htmlFor="code">Confirmation Code</label>
                <input
                    type="text"
                    id="code"
                    name="code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    required
                />
                <div className="button-row">
                    <button type="submit" className="confirm-button">Confirm</button>
                    <button type="button" className="resend-button" onClick={handleResend}>Resend</button>
                </div>
            </form>
        </div>
    );
};

export default ConfirmEmail;
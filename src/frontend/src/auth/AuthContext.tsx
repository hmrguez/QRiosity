import {createContext, ReactNode, useContext, useState} from 'react';
import AuthService from './AuthService';
import client from '../apolloClient';

interface AuthContextType {
	isAuthenticated: boolean;
	login: () => void;
	logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({children}: { children: ReactNode }) => {
	const authService = new AuthService(client);
	const [isAuthenticated, setIsAuthenticated] = useState(authService.isLoggedIn());

	const login = () => {
		setIsAuthenticated(true);
	};

	const logout = () => {
		authService.logout()
		setIsAuthenticated(false);
	};

	return (
		<AuthContext.Provider value={{isAuthenticated, login, logout}}>
			{children}
		</AuthContext.Provider>
	);
};

export const useAuth = () => {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return context;
};

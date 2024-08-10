import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import './index.css';
import 'primereact/resources/themes/saga-blue/theme.css'; // theme
import 'primereact/resources/primereact.min.css'; // core css
import 'primeicons/primeicons.css'; // icons
import 'primeflex/primeflex.css';
import {ApolloProvider} from '@apollo/client';
import client from "./apolloClient.ts";
import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Settings from "./settings/Settings.tsx";
import Login from "./auth/Login.tsx";
import Register from "./auth/Register.tsx";
import {AuthProvider} from './auth/AuthContext';
import PrivateRoute from './auth/PrivateRoute';

const router = createBrowserRouter([
	{
		path: "/",
		element: (
			<PrivateRoute>
				<App/>
			</PrivateRoute>
		),
		children: [
			{
				path: "setting",
				element: <Settings/>
			},
		]
	},
	{
		path: "login",
		element: <Login/>
	},
	{
		path: "register",
		element: <Register/>
	}
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
	<React.StrictMode>
		<ApolloProvider client={client}>
			<AuthProvider>
				<RouterProvider router={router}/>
			</AuthProvider>
		</ApolloProvider>
	</React.StrictMode>,
);
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
import Login from "./auth/Login.tsx";
import Register from "./auth/Register.tsx";
import {AuthProvider} from './auth/AuthContext';
import PrivateRoute from './auth/PrivateRoute';
import Roadmaps from "./learning/Roadmaps.tsx";
import MyLearning from "./learning/MyLearning.tsx";
import LandingPage from "./LandingPage.tsx";
import ConfirmEmail from "./auth/ConfirmEmail.tsx";
import Courses from "./learning/Courses.tsx";
import RoadmapDetail from "./learning/RoadmapDetail.tsx";
import RoadmapBuilder from "./learning/RoadmapBuilder.tsx";
import CourseCreator from "./learning/CourseCreator.tsx";
import Profile from "./auth/UserProfile.tsx";

const router = createBrowserRouter([
	{
		path: "/home",
		element: (
			<PrivateRoute>
				<App/>
			</PrivateRoute>
		),
		children: [
			{
				path: "profile",
				element: <Profile/>
			},
			{
				path: "my-learning",
				element: <MyLearning/>
			},
			{
				path: "roadmaps",
				element: <Roadmaps/>
			},
			{
				path: "courses",
				element: <Courses/>
			},
			{
				path: "roadmap/:id",
				element: <RoadmapDetail/>
			},
			{
				path: "build-roadmap",
				element: <RoadmapBuilder/>
			},
			{
				path: "create-course",
				element: <CourseCreator/>
			}
		]
	},
	{
		path: "/",
		element: <LandingPage/>
	},
	{
		path: "/login",
		element: <Login/>
	},
	{
		path: "/register",
		element: <Register/>
	},
	{
		path: "/confirm-email",
		element: <ConfirmEmail/>
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
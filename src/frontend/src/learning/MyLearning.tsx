import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService";
import AuthService from "../auth/AuthService";
import "./MyLearning.css";
import RoadmapList from "./RoadmapList.tsx";

const MyLearning = () => {
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client);

	const [roadmaps, setRoadmaps] = useState<any[]>([]);

	useEffect(() => {
		const fetchRoadmaps = async () => {
			const userId = authService.getCognitoUsername();
			try {
				const roadmaps = await learningService.getRoadmapsByUser(userId as string);
				setRoadmaps(roadmaps);
			} catch (error) {
				console.error(error);
			}
		};

		fetchRoadmaps();
	}, []);

	return (
		<div className="my-learning">
			<h1>My Roadmaps</h1>
			<RoadmapList roadmaps={roadmaps} myLearning={true}/>
		</div>
	);
};

export default MyLearning;
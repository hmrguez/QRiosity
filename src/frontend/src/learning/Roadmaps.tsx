import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import {Roadmap} from "./Roadmap.tsx";
import RoadmapList from "./RoadmapList.tsx";
import AuthService from "../auth/AuthService.tsx";
import "./Roadmaps.css"

const Roadmaps = () => {
	const client = useApolloClient();
	const authService = new AuthService(client)
	const learningService = new LearningService(client);

	const [roadmaps, setRoadmaps] = useState<Roadmap[]>([]);

	useEffect(() => {
		learningService.getRoadmapFeed(authService.getCognitoUsername() as string).then((roadmaps) => {
			setRoadmaps(roadmaps);
		}).catch((error) => {
			console.error(error);
		});
	}, []);


	return (
		// Use the new roadmap-list component
		<div className="roadmaps mt-5">
			<h1>Roadmaps</h1>
			<RoadmapList roadmaps={roadmaps} myLearning={false}/>
		</div>
	);
}

export default Roadmaps;
import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import "./Roadmaps.css";
import {Roadmap} from "./Roadmap.tsx";
import RoadmapList from "./RoadmapList.tsx";

const Roadmaps = () => {
	const client = useApolloClient();
	const learningService = new LearningService(client);

	const [roadmaps, setRoadmaps] = useState<Roadmap[]>([]);

	useEffect(() => {
		learningService.getRoadmaps().then((roadmaps) => {
			setRoadmaps(roadmaps);
		}).catch((error) => {
			console.error(error);
		});
	}, []);


	return (
		// Use the new roadmap-list component
		<div>
			<RoadmapList roadmaps={roadmaps} liked={[]}/>
		</div>
	);
}

export default Roadmaps;
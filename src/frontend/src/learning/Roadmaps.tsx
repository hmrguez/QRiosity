import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import {useNavigate} from "react-router-dom";
import LearningService from "./LearningService.tsx";
import "./Roadmaps.css";
import {Course} from "./Courses.tsx";

// Define the Roadmap interface
export interface Roadmap {
	id: string;
	title: string;
	author: string;
	courses: Course[];
	courseIDs: string[];
	topics: string[];
	isCustom: boolean;
	createdBy: string;
	likes: number;
	difficulty: string;
}

const Roadmaps = () => {
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const navigate = useNavigate();

	const [roadmaps, setRoadmaps] = useState<Roadmap[]>([]);

	useEffect(() => {
		learningService.getRoadmaps().then((roadmaps) => {
			setRoadmaps(roadmaps);
		}).catch((error) => {
			console.error(error);
		});
	}, []);

	const handleRoadmapClick = (id: string) => {
		navigate(`/home/roadmap/${id}`);
	};

	return (
		<div className="roadmap-grid">
			{roadmaps.map((roadmap) => (
				<div className="roadmap-card" key={roadmap.id} onClick={() => handleRoadmapClick(roadmap.id)}>
					<div className="roadmap-content">
						<h2 className="roadmap-title">{roadmap.title}</h2>
						<p className="roadmap-author">by {roadmap.author}</p>
						<div className="roadmap-topics">
							{roadmap.topics.map((topic, idx) => (
								<span className="topic-chip" key={idx}>{topic}</span>
							))}
						</div>
						<div className="roadmap-meta">
							<span className="roadmap-difficulty">{roadmap.difficulty}</span>
							<span className="roadmap-likes">{roadmap.likes}</span>
						</div>
						{roadmap.isCustom && <p className="roadmap-custom">Custom Roadmap</p>}
					</div>
				</div>
			))}
		</div>
	);
}

export default Roadmaps;
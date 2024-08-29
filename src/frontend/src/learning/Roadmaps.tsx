import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import {useNavigate} from "react-router-dom";
import {Button} from 'primereact/button';
import LearningService from "./LearningService.tsx";
import "./Roadmaps.css";
import AuthService from "../auth/AuthService.tsx";
import {Roadmap} from "./Roadmap.tsx";

const Roadmaps = () => {
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client)
	const navigate = useNavigate();

	const [roadmaps, setRoadmaps] = useState<Roadmap[]>([]);
	const [loading, setLoading] = useState<string | null>(null);

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

	const handleLikeClick = async (roadmapId: string) => {
		const userId = authService.getCognitoUsername();
		setLoading(roadmapId);
		try {
			const success = await learningService.userLikedRoadmap(userId as string, roadmapId);
			if (success) {
				setRoadmaps((prevRoadmaps) =>
					prevRoadmaps.map((roadmap) =>
						roadmap.id === roadmapId ? {...roadmap, likes: roadmap.likes + 1} : roadmap
					)
				);
			}
		} catch (error) {
			console.error(error);
		} finally {
			setLoading(null);
		}
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
						<Button
							icon="pi pi-thumbs-up"
							label="Like"
							className="p-button-rounded p-button-success p-button-text"
							onClick={(e) => {
								e.stopPropagation();
								handleLikeClick(roadmap.id);
							}}
							loading={loading === roadmap.id}
							disabled={loading === roadmap.id}
						/>
					</div>
				</div>
			))}
		</div>
	);
}

export default Roadmaps;
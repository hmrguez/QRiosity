import React, {useState} from "react";
import "./RoadmapList.css";
import {useNavigate} from "react-router-dom";
import {Roadmap} from "./Roadmap.tsx";
import {Button} from "primereact/button";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import AuthService from "../auth/AuthService.tsx";

interface RoadmapListProps {
	roadmaps: Roadmap[];
	myLearning: boolean;
}

const RoadmapList: React.FC<RoadmapListProps> = ({roadmaps, myLearning}) => {
	const navigate = useNavigate();
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client);
	const [loading, setLoading] = useState<string | null>(null);

	const handleLikeClick = async (roadmapId: string) => {
		const userId = authService.getCognitoUsername();
		setLoading(roadmapId);
		try {
			const success = await learningService.userLikedRoadmap(userId as string, roadmapId);
			if (success) {
				const roadmapIndex = roadmaps.findIndex((roadmap) => roadmap.id === roadmapId);
				roadmaps[roadmapIndex].likes += 1;
				roadmaps[roadmapIndex].liked = true;
			}
		} catch (error) {
			console.error(error);
		} finally {
			setLoading(null);
		}
	};

	const onRoadmapClick = (id: string) => {
		navigate(`/home/roadmap/${id}`);
	};

	return (
		<div className="roadmaps-grid">
			{roadmaps.map((roadmap) => (
				<div key={roadmap.id} className="roadmap-card" onClick={() => onRoadmapClick(roadmap.id)}>
					<img src={roadmap.imageUrl} alt="Roadmap Image" className="roadmap-image"/>
					<div className="roadmap-content">
						<h2 className="roadmap-title ">{roadmap.title}</h2>
						<div className="roadmap-author ">
							By {roadmap.author}
							{/*{roadmap.isVerified && <span className="verified-tag">✓</span>}*/}
						</div>
						<div className="roadmap-meta ">
							<span className="roadmap-difficulty">{roadmap.difficulty}</span>
							{roadmap.isCustom && <span className="custom-tag">Custom</span>}
						</div>
						<div className="roadmap-topics ">
							{roadmap.topics.slice(0, 5).map((topic, idx) => (
								<span key={idx} className="topic-chip">{topic}</span>
							))}
							{roadmap.topics.length > 5 && (
								<span className="topic-chip">+{roadmap.topics.length - 5} more</span>
							)}
						</div>
						<div className="roadmap-actions">
							{myLearning || roadmap.liked ? (
								<Button
									style={{color: '#FF6B6B'}}
									className="like-button"
									icon="pi pi-heart-fill"
								/>
							) : (
								<Button
									className="like-button"
									icon="pi pi-heart-fill"
									onClick={(e) => {
										e.stopPropagation();
										handleLikeClick(roadmap.id);
									}}
									loading={loading === roadmap.id}
									disabled={loading === roadmap.id}
								/>
							)}
							<span className="like-counter">{roadmap.likes}</span>
						</div>
					</div>
				</div>
			))}
		</div>
	);
};

export default RoadmapList;
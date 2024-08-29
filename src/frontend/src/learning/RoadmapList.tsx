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
	liked: boolean[];
}

const RoadmapList: React.FC<RoadmapListProps> = ({roadmaps, liked}) => {
	const navigate = useNavigate();
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client)
	const [loading, setLoading] = useState<string | null>(null);

	const handleLikeClick = async (roadmapId: string) => {
		const userId = authService.getCognitoUsername();
		setLoading(roadmapId);
		try {
			const success = await learningService.userLikedRoadmap(userId as string, roadmapId);
			if (success) {
				const roadmapIndex = roadmaps.findIndex((roadmap) => roadmap.id === roadmapId);
				roadmaps[roadmapIndex].likes += 1;
				liked[roadmapIndex] = true;
			}
		} catch (error) {
			console.error(error);
		} finally {
			setLoading(null);
		}
	};

	const onRoadmapClick = (id: string) => {
		navigate(`/home/roadmap/${id}`);
	}

	return (
		<div className="roadmap-list-wrapper">
			<div className="container">
				<h1>My Roadmaps</h1>
				<div className="roadmap-list">
					{roadmaps.map((roadmap, index) => (
						<div key={roadmap.id} className="roadmap-card" onClick={() => onRoadmapClick(roadmap.id)}>
							<div className="card-image">
								<img src={roadmap.imageUrl} alt="Roadmap Image"/>
							</div>
							<div className="card-content">
								<h2 className="card-title">{roadmap.title}</h2>
								<div className="card-meta">
									<span className="author">{roadmap.author}</span>
									{roadmap.isCustom && <span className="custom-badge">Custom</span>}
									<span className="difficulty">{roadmap.difficulty}</span>
								</div>
								<div className="card-footer">
									<div className="like-section">
										{liked[index] ?
											<Button
												style={{color: '#FF6B6B'}}
												className="like-button"
												icon="pi pi-heart-fill"
											/> :
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
										}
										<span className="like-counter">{roadmap.likes}</span>
									</div>
									<div className="topics">
										{roadmap.topics.map((topic, idx) => (
											<span key={idx} className="topic-tag">{topic}</span>
										))}
									</div>
								</div>
							</div>
						</div>
					))}
				</div>
			</div>
		</div>
	);
};

export default RoadmapList;
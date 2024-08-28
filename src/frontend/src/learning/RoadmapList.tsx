import React from "react";
import "./RoadmapList.css";
import {useNavigate} from "react-router-dom";

interface Roadmap {
	id: string;
	title: string;
	author: string;
	topics: string[];
	difficulty: string;
	likes: number;
	isCustom?: boolean;
}

interface RoadmapListProps {
	roadmaps: Roadmap[];
}

const RoadmapList: React.FC<RoadmapListProps> = ({roadmaps}) => {
	const navigate = useNavigate();

	const onRoadmapClick = (id: string) => {
		navigate(`/home/roadmap/${id}`);
	}

	return (
		<div className="container">
			<h1>My Roadmaps</h1>
			<div className="roadmap-list">
				{roadmaps.map((roadmap) => (
					<div key={roadmap.id} className="roadmap-card" onClick={() => onRoadmapClick(roadmap.id)}>
						<div className="roadmap-title">{roadmap.title}</div>
						<div className="roadmap-author">By: {roadmap.author}</div>
						<div className="roadmap-topics">
							{roadmap.topics.map((topic, idx) => (
								<span key={idx} className="topic">{topic}</span>
							))}
						</div>
						<div className="roadmap-meta">
							<span className="difficulty">{roadmap.difficulty}</span>
							<span className="likes">{roadmap.likes}</span>
							{roadmap.isCustom && <span className="custom-badge">Custom</span>}
						</div>
					</div>
				))}
			</div>
		</div>
	);
};

export default RoadmapList;
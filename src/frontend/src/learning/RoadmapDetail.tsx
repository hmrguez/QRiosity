import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import LearningService from './LearningService';
import './RoadmapDetail.css';
import {useApolloClient} from "@apollo/client";

import {Roadmap} from "./Roadmap.tsx";
import OutOfViews from "../utils/OutOfViews.tsx";
import {BlockUI} from "primereact/blockui";
import Loading from "../utils/Loading.tsx";

enum RoadmapStatus {
	Loading,
	NoViewsRemaining,
	Fetched,
}

const RoadmapDetail = () => {
	const {id} = useParams();
	const [roadmap, setRoadmap] = useState<Roadmap | null>(null);
	const [status, setStatus] = useState<RoadmapStatus>(RoadmapStatus.Loading)
	const apolloClient = useApolloClient();
	const learningService = new LearningService(apolloClient);

	useEffect(() => {
		const fetchRoadmap = async () => {

			if (!id) {
				return;
			}

			const roadmapData = await learningService.getRoadmapById(id);

			if (roadmapData === "user has no views remaining") {
				setStatus(RoadmapStatus.NoViewsRemaining)
			} else if (roadmapData === "roadmap not found") {
				// Navigate to not found page
			} else if (roadmapData instanceof String) {
				// Show toast
			}

			setRoadmap(roadmapData as Roadmap);
			setStatus(RoadmapStatus.Fetched)
		};

		fetchRoadmap();
	}, [id]);

	const handleCourseClick = (courseUrl: string) => {
		window.open(courseUrl, '_blank');
	};

	if (!roadmap) {
		return (<Loading/>)
	}

	return (
		<div>
			<BlockUI blocked={status !== RoadmapStatus.Fetched}>
				<div className="roadmap-container">
					<div className="roadmap-courses">
						<ul className="course-list">
							{roadmap.courses.map(course => (
								<li key={course.id} className="course-item"
									onClick={() => handleCourseClick(course.url)}>
									<h3 className="course-title">{course.title}</h3>
									<p className="course-description">{course.description}</p>
									<div className="course-meta">
										<span className="course-difficulty">{course.difficulty}</span>
										<span className="course-duration">{course.duration}</span>
									</div>
								</li>
							))}
						</ul>
					</div>
					<div className="roadmap-details">
						<h2 className="detail-roadmap-title">{roadmap.title}</h2>
						<p className="roadmap-author">by {roadmap.author}</p>
						<div className="roadmap-stats">
							<div className="stat-item">
								<div className="stat-value">{roadmap.courses.length}</div>
								<div className="stat-label">Courses</div>
							</div>
							<div className="stat-item">
								<div className="stat-value">{roadmap.difficulty}</div>
								<div className="stat-label">Difficulty</div>
							</div>
							<div className="stat-item">
								<div className="stat-value">{roadmap.likes}</div>
								<div className="stat-label">Likes</div>
							</div>
						</div>
						<div className="roadmap-topics">
							{roadmap.topics.map(topic => (
								<span key={topic} className="topic-chip">{topic}</span>
							))}
						</div>
						{roadmap.isCustom && (<p className="roadmap-custom">Custom</p>)}
						<i className="detail-roadmap-description">{roadmap.description}</i>
					</div>
				</div>
			</BlockUI>
			{status === RoadmapStatus.NoViewsRemaining && (<OutOfViews/>)}
		</div>
	);
};

export default RoadmapDetail;
import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import LearningService from './LearningService';
import './RoadmapDetail.css';
import {useApolloClient} from "@apollo/client";

import {Roadmap} from "./Roadmap.tsx";

const RoadmapDetail = () => {
	const {id} = useParams();
	const [roadmap, setRoadmap] = useState<Roadmap | null>(null);

	const apolloClient = useApolloClient();
	const learningService = new LearningService(apolloClient);

	useEffect(() => {
		const fetchRoadmap = async () => {

			if (!id) {
				return;
			}

			const roadmapData = await learningService.getRoadmapById(id);
			setRoadmap(roadmapData);
		};

		fetchRoadmap();
	}, [id]);

	if (!roadmap) {
		return <div>Loading...</div>;
	}

	const handleCourseClick = (courseUrl: string) => {
		window.open(courseUrl, '_blank');
	};

	return (
		<div className="roadmap-container">
			<div className="roadmap-courses">
				<ul className="course-list">
					{roadmap.courses.map(course => (
						<li key={course.id} className="course-item" onClick={() => handleCourseClick(course.url)}>
							<h3 className="course-title">{course.title}</h3>
							<p className="course-description">{course.description}</p>
							<div className="course-meta">
								<span className="course-difficulty">{course.difficulty}</span>
								<span className="course-duration">{course.duration} hours</span>
							</div>
						</li>
					))}
				</ul>
			</div>
			<div className="roadmap-details">
				<h2 className="roadmap-title">{roadmap.title}</h2>
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
				<p className="roadmap-custom">{roadmap.isCustom ? 'Custom Roadmap' : 'Standard Roadmap'}</p>
			</div>
		</div>
	);
};

export default RoadmapDetail;
import {useEffect, useRef, useState} from 'react';
import {useParams} from 'react-router-dom';
import LearningService from './LearningService';
import './RoadmapDetail.css';
import {useApolloClient} from "@apollo/client";

import {Roadmap} from "./Roadmap.tsx";
import OutOfViews from "../utils/OutOfViews.tsx";
import Loading from "../utils/Loading.tsx";
import {Toast} from "primereact/toast";

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

	const toast = useRef<Toast>(null);

	useEffect(() => {
		const fetchRoadmap = async () => {

			if (!id) {
				window.location.href = '/404';
				return
			}

			try {
				const roadmapData = await learningService.getRoadmapById(id);
				setRoadmap(roadmapData as Roadmap);
				setStatus(RoadmapStatus.Fetched)
			} catch (e: any) {
				if (e.message == "user has no views remaining") {
					setStatus(RoadmapStatus.NoViewsRemaining)
				} else if (e.message == "roadmap not found") {
					window.location.href = '/404';
				} else {
					toast.current?.show({severity: 'error', summary: 'Error', detail: e.message, life: 3000});
				}
			}
		};

		fetchRoadmap();
	}, [id]);

	const handleCourseClick = (courseUrl: string) => {
		window.open(courseUrl, '_blank');
	};

	if (!roadmap && status === RoadmapStatus.Loading) {
		return (<Loading/>)
	}

	return (
		<div>
			<Toast ref={toast}/>
			{status === RoadmapStatus.Fetched && (
				<div className="roadmap-container">
					<div className="roadmap-courses">
						<ul className="course-list">
							{roadmap?.courses.map(course => (
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
						<h2 className="detail-roadmap-title">{roadmap?.title}</h2>
						<p className="roadmap-author">by {roadmap?.author}</p>
						<div className="roadmap-stats">
							<div className="stat-item">
								<div className="stat-value">{roadmap?.courses.length}</div>
								<div className="stat-label">Courses</div>
							</div>
							<div className="stat-item">
								<div className="stat-value">{roadmap?.difficulty}</div>
								<div className="stat-label">Difficulty</div>
							</div>
							<div className="stat-item">
								<div className="stat-value">{roadmap?.likes}</div>
								<div className="stat-label">Likes</div>
							</div>
						</div>
						<div className="roadmap-topics">
							{roadmap?.topics.map(topic => (
								<span key={topic} className="topic-chip">{topic}</span>
							))}
						</div>
						{roadmap?.isCustom && (<p className="roadmap-custom">Custom</p>)}
						<i className="detail-roadmap-description">{roadmap?.description}</i>
					</div>
				</div>
			)}
			{status === RoadmapStatus.NoViewsRemaining && (<OutOfViews/>)}
		</div>
	);
};

export default RoadmapDetail;
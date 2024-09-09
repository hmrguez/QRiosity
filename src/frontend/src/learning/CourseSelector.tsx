// CourseSelector.tsx
import React, {useEffect, useState} from 'react';
import {Course} from './Course.tsx';
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import AuthService from "../auth/AuthService.tsx";
import './CourseSelector.css';

interface CourseSelectorProps {
	isOpen: boolean;
	onClose: () => void;
	onCourseSelect: (course: Course) => void;
}

const CourseSelector: React.FC<CourseSelectorProps> = ({isOpen, onClose, onCourseSelect}) => {

	const apolloClient = useApolloClient();
	const learningService = new LearningService(apolloClient);
	const authService = new AuthService(apolloClient);

	const [courses, setCourses] = useState<Course[]>([]);
	const [pagination, setPagination] = useState({page: 0, perPage: 12, lastEvaluatedKey: null as string | null});

	const fetchCourses = () => {
		const cognitoUsername = authService.getCognitoUsername() as string;
		learningService.getCourses(cognitoUsername, pagination).then((result) => {
			setCourses(result.courses);
			setPagination(result.pagination);
		}).catch((error) => {
			console.error(error);
		});
	};

	useEffect(() => {
		fetchCourses();
	}, []);

	const handleNextPage = () => {
		setPagination(prev => ({...prev, page: prev.page + 1}));
		console.log(pagination);
		fetchCourses();
	};

	const handlePreviousPage = () => {
		setPagination(prev => ({...prev, page: Math.max(prev.page - 1, 0)}));
		fetchCourses();
	};

	if (!isOpen) return null;

	return (
		<div id="course-selector" className="course-selector" onClick={(e) => {
			if (e.target === e.currentTarget) onClose();
		}}>
			<div className="course-selector-content">
				<h2>Select a Course</h2>
				{/*<input type="text" id="course-search" placeholder="Search courses..." onChange={handleSearch}/>*/}
				<ul id="course-list" className="selectable-course-list">
					{courses.map(course => (
						<li key={course.id} className="selectable-course-item" onClick={() => onCourseSelect(course)}>
							<h3>{course.title}</h3>
							<p>{course.description}</p>
							<div className="course-meta">
								<span>{course.difficulty}</span>
								<span>{course.duration} hours</span>
							</div>
						</li>
					))}
				</ul>

				<div className="selector-paginator">
					<button onClick={handlePreviousPage}>&laquo;</button>
					<button onClick={handleNextPage}>&raquo;</button>
				</div>

				<button id="close-selector" className="close-button" onClick={onClose}>Close</button>
			</div>
		</div>
	);
};

export default CourseSelector;
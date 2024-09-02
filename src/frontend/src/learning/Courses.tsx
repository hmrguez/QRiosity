import './Courses.css';
import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import AuthService from "../auth/AuthService.tsx";

export interface Course {
	id: string;
	title: string;
	description: string;
	topics: string[];
	difficulty: string;
	duration: string;
	language: string;
	isFree: boolean;
	url: string;
}

const Courses = () => {
	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client);

	const [courses, setCourses] = useState<Course[]>([]);
	const [pagination, setPagination] = useState({page: 0, perPage: 3, lastEvaluatedKey: null as string | null});

	const fetchCourses = () => {
		const cognitoUsername = authService.getCognitoUsername() as string;
		learningService.getCourses(cognitoUsername, pagination).then((result) => {
			console.log(result);
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
		fetchCourses();
	};

	const handlePreviousPage = () => {
		setPagination(prev => ({...prev, page: Math.max(prev.page - 1, 0)}));
		fetchCourses();
	};

	const handleCardClick = (url: string) => {
		window.open(url, '_blank');
	};

	return (
		<div>
			<div className="course-grid">
				{courses.map((course, index) => (
					<div className="course-card" key={index} onClick={() => handleCardClick(course.url)}>
						<div className="course-image"
							 style={{backgroundImage: `url("https://via.placeholder.com/40")`}}></div>
						<div className="course-content">
							<h2 className="course-title">{course.title}</h2>
							<p className="course-description">{course.description}</p>
							<div className="course-topics">
								{course.topics.map((topic, idx) => (
									<span className="topic-chip" key={idx}>{topic}</span>
								))}
							</div>
							<div className="course-meta">
								<span className="course-difficulty">{course.difficulty}</span>
								<span className="course-duration">{course.duration}</span>
							</div>
							<div className="course-meta">
								<span className="course-language">{course.language}</span>
								<span className="course-price">{course.isFree ? "Free" : ""}</span>
							</div>
						</div>
					</div>
				))}
			</div>

			<div className="paginator">
				<button onClick={handlePreviousPage}>&laquo;</button>
				<button className="active">1</button>
				<button>2</button>
				<button>3</button>
				<button onClick={handleNextPage}>&raquo;</button>
			</div>
		</div>
	);
}

export default Courses;
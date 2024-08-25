import './Courses.css';
import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";

export interface Course {
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

    const [courses, setCourses] = useState<Course[]>([]);

    useEffect(() => {
        learningService.getCourses().then((courses) => {
            setCourses(courses);
        }).catch((error) => {
            console.error(error);
        });
    }, []);

    const handleCardClick = (url: string) => {
        window.open(url, '_blank');
    };

    return (
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
    );
}

export default Courses;
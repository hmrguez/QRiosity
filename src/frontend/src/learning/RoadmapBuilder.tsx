import {useEffect, useRef, useState} from 'react';
import './RoadmapBuilder.css';
import {useApolloClient} from "@apollo/client";
import LearningService from "./LearningService.tsx";
import {Button} from "primereact/button";
import AuthService from "../auth/AuthService.tsx";
import {Toast} from "primereact/toast";
import {useNavigate} from "react-router-dom";
import FileUpload from "../utils/FileInput.tsx";
import CreatorService from "../creator/CreatorService.tsx";
import {Course} from "./Course.tsx";

const RoadmapBuilder = () => {
	const fileUploadRef = useRef<{ handleSubmit: () => void }>(null);


	const [roadmapTitle, setRoadmapTitle] = useState('');
	const [roadmapDescription, setRoadmapDescription] = useState('');
	const [roadmapDifficulty, setRoadmapDifficulty] = useState('beginner');
	const [roadmapTopics, setRoadmapTopics] = useState('');
	const [isModalOpen, setIsModalOpen] = useState(false);
	const [courses, setCourses] = useState<Course[]>([]);
	const [filteredCourses, setFilteredCourses] = useState(courses);
	const [addedCourses, setAddedCourses] = useState<Course[]>([]);
	const [loading, setLoading] = useState(false);

	const openModal = () => setIsModalOpen(true);
	const closeModal = () => setIsModalOpen(false);

	const apolloClient = useApolloClient();
	const learningService = new LearningService(apolloClient);
	const authService = new AuthService(apolloClient);
	const creatorService = new CreatorService(apolloClient)

	const [prompt, setPrompt] = useState('');
	const [usagesRemaining, setUsagesRemaining] = useState(0);
	const [totalUsages, setTotalUsages] = useState(20);
	const [loadingAutoGenerate, setLoadingAutoGenerate] = useState(false);

	const toast = useRef<Toast>(null);

	const handleAutoGenerate = async () => {
		setLoadingAutoGenerate(true);
		try {
			const userId = authService.getCognitoUsername();
			const data = await learningService.customRoadmapRequested(prompt, userId as string);
			setRoadmapTitle(data.title);
			setAddedCourses(data.courses);
			setRoadmapTopics(data.topics.join(', '));
			setRoadmapDifficulty(data.difficulty);
			setRoadmapDescription(data.description);
			setUsagesRemaining(usagesRemaining - 1);
		} catch (error) {
			console.error(error);
			toast.current?.show({
				severity: 'error',
				summary: 'Error',
				detail: 'Failed to auto-generate roadmap',
				life: 3000
			});
		} finally {
			setLoadingAutoGenerate(false);
		}
	};

	const addCourseToRoadmap = (course: Course) => {
		setAddedCourses([...addedCourses, course]);
		closeModal();
	};

	const removeCourseFromRoadmap = (courseId: string) => {
		setAddedCourses(addedCourses.filter(course => course.id !== courseId));
	};

	const handleSearch = (e: any) => {
		const searchTerm = e.target.value.toLowerCase();
		const filtered = courses.filter(course =>
			course.title.toLowerCase().includes(searchTerm) ||
			course.description.toLowerCase().includes(searchTerm)
		);
		setFilteredCourses(filtered);
	};

	const router = useNavigate();

	const saveRoadmap = async () => {
		setLoading(true);
		let roadmapInput = {
			id: crypto.randomUUID(),
			title: roadmapTitle,
			author: authService.getUsername(),
			authorId: authService.getCognitoUsername(),
			courseIDs: addedCourses.map(course => course.id),
			topics: roadmapTopics.split(',').map(topic => topic.trim()),
			isCustom: true,
			createdBy: authService.getUsername(),
			likes: 0,
			difficulty: roadmapDifficulty,
			imageUrl: '',
			description: roadmapDescription
		};


		try {
			const filename = await fileUploadRef.current?.handleSubmit();

			if (filename) {
				roadmapInput.imageUrl = filename;
			}

			await learningService.upsertRoadmap(roadmapInput);
			toast.current?.show({severity: 'success', summary: 'Success', detail: 'Roadmap saved', life: 3000});
			router(`/home/roadmap/${roadmapInput.id}`);
		} catch (error) {
			console.error(error);
			toast.current?.show({severity: 'error', summary: 'Error', detail: 'Failed to save roadmap', life: 3000});
		} finally {
			setLoading(false);
		}
	};

	useEffect(() => {
		const fetchCourses = async () => {

			const cognitoUsername = authService.getCognitoUsername();
			const pagination = {page: 0, perPage: 10};

			const fetchedCourses = await learningService.getCourses(cognitoUsername as string, pagination);
			setCourses(fetchedCourses.courses);
			setFilteredCourses(fetchedCourses.courses);
		};

		const fetchGenUsage = async () => {
			const {role, usages} = await creatorService.getGenUsagesRemaining();
			setUsagesRemaining(usages);
			setTotalUsages(role == 1 ? 3 : 6);
		}

		fetchGenUsage();

		fetchCourses();
	}, []);

	const moveCourseUp = (index: number) => {
		if (index > 0) {
			const newCourses = [...addedCourses];
			[newCourses[index - 1], newCourses[index]] = [newCourses[index], newCourses[index - 1]];
			setAddedCourses(newCourses);
		}
	};

	const moveCourseDown = (index: number) => {
		if (index < addedCourses.length - 1) {
			const newCourses = [...addedCourses];
			[newCourses[index + 1], newCourses[index]] = [newCourses[index], newCourses[index + 1]];
			setAddedCourses(newCourses);
		}
	};

	return (
		<div className="roadmap-container">
			<Toast ref={toast}/>
			<div className="roadmap-courses">
				<ul className="course-list">
					{addedCourses.map((course, index) => (
						<li key={course.id} className="course-item">
							<div className="course-content">
								<h3 className="course-title">{course.title}</h3>
								<p className="course-description">{course.description}</p>
								<div className="course-meta">
									<span className="course-difficulty">{course.difficulty}</span>
									<span className="course-duration">{course.duration}</span>
								</div>
							</div>
							<div className="course-actions">
								<Button className="p-button-icon p-button-contrast" icon="pi pi-chevron-up"
										onClick={() => moveCourseUp(index)}/>
								<Button className="p-button-icon p-button-danger" icon="pi pi-times"
										onClick={() => removeCourseFromRoadmap(course.id)}/>
								<Button className="p-button-icon p-button-contrast" icon="pi pi-chevron-down"
										onClick={() => moveCourseDown(index)}/>
							</div>
						</li>
					))}
				</ul>
				<div className="add-course" onClick={openModal}>+ Add new course</div>
			</div>
			<div className="roadmap-details">
				<h3>Generate with AI</h3>
				<div className="prompt">
					<input type="text"
						   id="prompt"
						   name="prompt"
						   value={prompt}
						   onChange={(e) => setPrompt(e.target.value)}
						   placeholder="Roadmap Topic"
					/>
					<Button className="auto-generate-btn"
							label="Auto Generate"
							onClick={handleAutoGenerate}
							loading={loadingAutoGenerate}
							disabled={usagesRemaining == 0 || loadingAutoGenerate}>
					</Button>
				</div>
				<div className="progress-container">
					<div className="progress-bar">
						<div
							className="progress-fill"
							style={{width: `${(usagesRemaining / totalUsages) * 100}%`}}
						></div>
					</div>
					<div className="progress-label">
						<span>Uses Left</span>
						<span>{usagesRemaining}</span>
					</div>
				</div>

				<br/>

				<h3>Roadmap Details</h3>
				<div className="input-group">
					<label htmlFor="roadmap-title">Roadmap Title</label>
					<input
						type="text"
						id="roadmap-title"
						name="roadmap-title"
						placeholder="Enter roadmap title"
						value={roadmapTitle}
						onChange={(e) => setRoadmapTitle(e.target.value)}
					/>
				</div>
				<div className="input-group">
					<label htmlFor="roadmap-description">Description</label>
					<textarea
						id="roadmap-description"
						name="roadmap-description"
						placeholder="Enter roadmap description"
						value={roadmapDescription}
						onChange={(e) => setRoadmapDescription(e.target.value)}
					></textarea>
				</div>
				<div className="input-group">
					<label htmlFor="roadmap-difficulty">Difficulty</label>
					<select
						id="roadmap-difficulty"
						name="roadmap-difficulty"
						value={roadmapDifficulty}
						onChange={(e) => setRoadmapDifficulty(e.target.value)}
					>
						<option value="Beginner">Beginner</option>
						<option value="Intermediate">Intermediate</option>
						<option value="Advanced">Advanced</option>
					</select>
				</div>
				<div className="input-group">
					<label htmlFor="roadmap-topics">Topics (comma-separated)</label>
					<input
						type="text"
						id="roadmap-topics"
						name="roadmap-topics"
						placeholder="e.g. JavaScript, Web Development, Frontend"
						value={roadmapTopics}
						onChange={(e) => setRoadmapTopics(e.target.value)}
					/>
				</div>
				<FileUpload ref={fileUploadRef}/>

				<Button severity="contrast" onClick={saveRoadmap} disabled={loading} loading={loading}>
					{loading ? 'Saving...' : 'Save Roadmap'}
				</Button>
			</div>

			{isModalOpen && (
				<div id="course-selector" className="course-selector" onClick={(e) => {
					if (e.target === e.currentTarget) closeModal();
				}}>
					<div className="course-selector-content">
						<h2>Select a Course</h2>
						<input type="text" id="course-search" placeholder="Search courses..."
							   onChange={handleSearch}/>
						<ul id="course-list" className="selectable-course-list">
							{filteredCourses.map(course => (
								<li key={course.id} className="selectable-course-item"
									onClick={() => addCourseToRoadmap(course)}>
									<h3>{course.title}</h3>
									<p>{course.description}</p>
									<div className="course-meta">
										<span>{course.difficulty}</span>
										<span>{course.duration} hours</span>
									</div>
								</li>
							))}
						</ul>
						<button id="close-selector" className="close-button" onClick={closeModal}>Close</button>
					</div>
				</div>
			)}
		</div>
	);
};

export default RoadmapBuilder;
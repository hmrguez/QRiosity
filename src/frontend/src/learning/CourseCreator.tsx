import {useRef, useState} from 'react';
import './CourseCreator.css';
import LearningService from '../learning/LearningService';
import {useApolloClient} from "@apollo/client";
import AuthService from "../auth/AuthService.tsx";
import {Button} from "primereact/button";
import {Toast} from "primereact/toast";

const CourseCreator = () => {
	const [title, setTitle] = useState('');
	const [description, setDescription] = useState('');
	const [topics, setTopics] = useState('');
	const [difficulty, setDifficulty] = useState('');
	const [duration, setDuration] = useState('');
	const [language, setLanguage] = useState('');
	const [isFree, setIsFree] = useState(false);
	const [url, setUrl] = useState('');
	const [source, setSource] = useState('');
	const [loading, setLoading] = useState(false);

	const client = useApolloClient();
	const learningService = new LearningService(client);
	const authService = new AuthService(client);

	const toast = useRef<Toast>(null);

	const handleSubmit = async (e: any) => {
		e.preventDefault();
		const courseInput = {
			id: crypto.randomUUID(), // Generate or fetch the ID as needed
			title,
			description,
			topics: topics.split(',').map(topic => topic.trim()),
			difficulty,
			duration: parseInt(duration, 10),
			language,
			isFree,
			url,
			source: source, // Assuming source is 'self' for user-created courses
			author: authService.getUsername(), // Replace with actual current user
		};

		setLoading(true);

		try {
			const result = await learningService.upsertCourse(courseInput);
			if (result) {
				toast.current?.show({
					severity: 'success',
					summary: 'Success',
					detail: 'Course added successfully',
					life: 3000
				});
				setTitle('');
				setDescription('');
				setTopics('');
				setDifficulty('');
				setDuration('');
				setLanguage('');
				setIsFree(false);
				setUrl('');
				setSource('');
			}
		} catch (error) {
			console.error(error);
			toast.current?.show({
				severity: 'error',
				summary: 'Error',
				detail: 'Failed to add course',
				life: 3000
			});
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="course-form-container">
			<Toast ref={toast}/>
			<h2 className="form-title">Add New Course</h2>
			<form id="course-form" onSubmit={handleSubmit}>
				<div className="input-group">
					<label htmlFor="title">Title</label>
					<input
						type="text"
						id="title"
						name="title"
						value={title}
						onChange={(e) => setTitle(e.target.value)}
						required
					/>
				</div>
				<div className="input-group">
					<label htmlFor="url">URL</label>
					<input
						type="text"
						id="url"
						name="url"
						value={url}
						onChange={(e) => setUrl(e.target.value)}
						required
					/>
				</div>
				<div className="input-group">
					<label htmlFor="source">Source</label>
					<input
						type="text"
						id="source"
						name="source"
						placeholder="e.g. Udemy | FreeCodeCamp | YouTube"
						value={source}
						onChange={(e) => setSource(e.target.value)}
						required
					/>
				</div>
				<div className="input-group">
					<label htmlFor="description">Description</label>
					<textarea
						id="description"
						name="description"
						value={description}
						onChange={(e) => setDescription(e.target.value)}
						required
					></textarea>
				</div>
				<div className="input-group">
					<label htmlFor="topics">Topics (comma-separated)</label>
					<input
						type="text"
						id="topics"
						name="topics"
						value={topics}
						onChange={(e) => setTopics(e.target.value)}
						placeholder="e.g. JavaScript, Web Development, React"
					/>
				</div>
				<div className="input-group">
					<label htmlFor="difficulty">Difficulty</label>
					<select
						id="difficulty"
						name="difficulty"
						value={difficulty}
						onChange={(e) => setDifficulty(e.target.value)}
						required
					>
						<option value="">Select difficulty</option>
						<option value="Beginner">Beginner</option>
						<option value="Intermediate">Intermediate</option>
						<option value="Advanced">Advanced</option>
					</select>
				</div>
				<div className="input-group">
					<label htmlFor="duration">Duration (in hours)</label>
					<input
						type="number"
						id="duration"
						name="duration"
						value={duration}
						onChange={(e) => setDuration(e.target.value)}
						min="1"
						required
					/>
				</div>
				<div className="input-group">
					<label htmlFor="language">Language</label>
					<input
						type="text"
						id="language"
						name="language"
						value={language}
						onChange={(e) => setLanguage(e.target.value)}
						required
					/>
				</div>
				<div className="input-group">
					<label className="checkbox-label">
						<input
							type="checkbox"
							id="isFree"
							name="isFree"
							checked={isFree}
							onChange={(e) => setIsFree(e.target.checked)}
						/>
						Is this course free?
					</label>
				</div>
				<Button type="submit" className="submit-button" loading={loading} disabled={loading}>Add Course</Button>
			</form>
		</div>
	);
};

export default CourseCreator;
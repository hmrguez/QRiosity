// src/frontend/src/components/UserProfile.tsx
import React, {useState} from 'react';
import './UserProfile.css';
import {Button} from "primereact/button";

const UserProfile: React.FC = () => {
	const [topics, setTopics] = useState<string[]>(['JavaScript', 'Python', 'Machine Learning', 'Web Development']);

	const renderTopics = () => {
		return topics.map((topic, index) => (
			<div key={index} className="topic" onClick={() => editTopic(index)}>
				{topic}
			</div>
		));
	};

	const editTopic = (index: number) => {
		const newTopics = [...topics];
		const currentTopic = newTopics[index];
		newTopics[index] = (
			<input
				type="text"
				className="topic-input"
				defaultValue={currentTopic}
				onBlur={(e) => handleBlur(e, index)}
				onKeyPress={(e) => handleKeyPress(e, index)}
				autoFocus
			/>
		) as unknown as string;
		setTopics(newTopics);
	};

	const handleBlur = (e: React.FocusEvent<HTMLInputElement>, index: number) => {
		const newTopics = [...topics];
		if (e.target.value.trim()) {
			newTopics[index] = e.target.value.trim();
		} else {
			newTopics.splice(index, 1);
		}
		setTopics(newTopics);
	};

	const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>, index: number) => {
		if (e.key === 'Enter') {
			(e.target as HTMLInputElement).blur();
		}
	};

	const addTopic = () => {
		setTopics([...topics, 'New Topic']);
		setTimeout(() => editTopic(topics.length), 0);
	};

	return (
		<div className="profile-container">
			<div className="profile-header">
				<img src="/api/placeholder/100/100" alt="Profile Picture" className="profile-pic"/>
				<div className="profile-info">
					<h1>John Doe</h1>
					<h2>Software Developer</h2>
				</div>
			</div>
			<div className="topics-container">
				<h3 className="topics-title">Learning Topics</h3>
				<div className="topics" id="topics-list">
					{renderTopics()}
				</div>
				<div className="flex justify-content-between">
					<Button className="p-button-contrast p-button-rounded mt-3" label="Add Topic" onClick={addTopic}/>
					<Button className="p-button-contrast p-button-rounded mt-3" label="Save"/>
				</div>

			</div>
		</div>
	);
};

export default UserProfile;
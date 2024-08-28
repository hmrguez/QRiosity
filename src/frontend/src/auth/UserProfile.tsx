import React, {useEffect, useState} from 'react';
import './UserProfile.css';
import {Button} from "primereact/button";

const UserProfile: React.FC = () => {
	const [name, setName] = useState<string>('John Doe');
	const [isEditingName, setIsEditingName] = useState<boolean>(false);
	const [topics, setTopics] = useState<string[]>(['JavaScript', 'Python', 'Machine Learning', 'Web Development']);
	const [isChanged, setIsChanged] = useState<boolean>(false);

	useEffect(() => {
		setIsChanged(true);
	}, [name, topics]);

	const renderTopics = () => {
		return topics.map((topic, index) => (
			<div key={index} className="topic">
				{topic}
				<button className="delete-topic" onClick={() => deleteTopic(index)}>x</button>
			</div>
		));
	};

	const deleteTopic = (index: number) => {
		const newTopics = topics.filter((_, i) => i !== index);
		setTopics(newTopics);
	};

	const addTopic = () => {
		setTopics([...topics, 'New Topic']);
	};

	const saveChanges = () => {
		setIsEditingName(false);
		setIsChanged(false);
	};

	return (
		<div className="profile-container">
			<div className="profile-header">
				<img src="/api/placeholder/100/100" alt="Profile Picture" className="profile-pic"/>
				<div className="profile-info">
					{isEditingName ? (
						<input
							type="text"
							className="name-input"
							value={name}
							onChange={(e) => setName(e.target.value)}
							onBlur={() => setIsEditingName(false)}
							autoFocus
						/>
					) : (
						<h1 onClick={() => setIsEditingName(true)}>{name}</h1>
					)}
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
					<Button className="p-button-contrast p-button-rounded mt-3" label="Save" onClick={saveChanges}
							disabled={!isChanged}/>
				</div>
			</div>
		</div>
	);
};

export default UserProfile;
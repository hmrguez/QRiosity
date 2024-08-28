import {useParams} from "react-router-dom";
import {useEffect, useState} from "react";
import {useApolloClient} from "@apollo/client";
import AuthService from "./AuthService.tsx";
import {Button} from "primereact/button";
import "./UserProfile.css"

const UserProfile: React.FC = () => {
	const {name} = useParams<{ name: string }>();
	const [userName, setUserName] = useState<string>('');
	const [role, setRole] = useState<string>('');
	const [topics, setTopics] = useState<string[]>([]);
	const [isEditingName, setIsEditingName] = useState<boolean>(false);
	const [isChanged, setIsChanged] = useState<boolean>(false);
	const [loading, setLoading] = useState<boolean>(false);
	const [isEditable, setIsEditable] = useState<boolean>(false);

	const apolloClient = useApolloClient();
	const authService = new AuthService(apolloClient)

	useEffect(() => {
		const fetchProfile = async () => {
			if (!name) {
				// Route to root if no name is provided
				return;
			}

			const profile = await authService.getProfile(name);
			setUserName(profile.username);
			setTopics(profile.topics);

			switch (profile.role) {
				case 0:
					setRole('Rookie');
					break;
				case 1:
					setRole('Aspirant');
					break;
				case 2:
					setRole('Creator');
					break;
				default:
					setRole('Unknown');
			}

			const cognitoUsername = await authService.getCognitoUsername();
			setIsEditable(cognitoUsername === name);
		};
		fetchProfile();
	}, [name]);

	useEffect(() => {
		setIsChanged(true);
	}, [userName, topics, role]);

	const editTopic = (index: number, newValue: string) => {
		const newTopics = [...topics];
		newTopics[index] = newValue;
		setTopics(newTopics);
	};

	const handleBlur = (index: number) => {
		const newTopics = [...topics];
		if (!newTopics[index].trim()) {
			newTopics.splice(index, 1);
		}
		setTopics(newTopics);
	};

	const renderTopics = () => {
		return topics.map((topic, index) => (
			<div key={index} className="topic">
				{isEditable ? (
					<input
						type="text"
						value={topic}
						onChange={(e) => editTopic(index, e.target.value)}
						onBlur={() => handleBlur(index)}
						className="bg-transparent outline-none"
					/>
				) : (
					<span>{topic}</span>
				)}
				{isEditable && (
					<button className="delete-topic" onClick={() => deleteTopic(index)}>x</button>
				)}
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

	const saveChanges = async () => {
		setLoading(true);
		// Add logic to save changes here, e.g., call updateUser mutation
		await authService.updateUser({
			name: name as string,
			username: userName,
			role: role === 'Rookie' ? 0 : role === 'Aspirant' ? 1 : 2,
			topics: topics, // Set this based on your logic
		});
		setIsEditingName(false);
		setIsChanged(false);
		setLoading(false);
	};

	return (
		<div className="super-container">
			<div className="profile-container">
				<div className="profile-header">
					<img src="/api/placeholder/100/100" alt="Profile Picture" className="profile-pic"/>
					<div className="profile-info">
						{isEditingName ? (
							<input
								type="text"
								className="name-input"
								value={userName}
								onChange={(e) => setUserName(e.target.value)}
								onBlur={() => setIsEditingName(false)}
								autoFocus
							/>
						) : isEditable ? (
							<h1 onClick={() => setIsEditingName(true)}>{userName}</h1>
						) : (
							<h1>{userName}</h1>
						)}
						<h2>{role}</h2>
					</div>
				</div>
				<div className="topics-container">
					<h3 className="topics-title">Learning Topics</h3>
					<div className="topics" id="topics-list">
						{renderTopics()}
					</div>
					{isEditable && (
						<div className="flex justify-content-between">
							<Button className="p-button-contrast p-button-rounded mt-3" label="Add Topic"
									onClick={addTopic}/>
							<Button className="p-button-contrast p-button-rounded mt-3"
									label={loading ? "Saving..." : "Save"} onClick={saveChanges}
									disabled={!isChanged || loading}/>
						</div>
					)}
				</div>
			</div>
		</div>
	);
};

export default UserProfile;
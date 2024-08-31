import {SetStateAction, useEffect, useRef, useState} from 'react';
import {Link} from 'react-router-dom';
import './Navbar.css';
import DailyChallengeModal from "../learning/DailyChallengeModal.tsx";
import {useAuth} from "../auth/AuthContext.tsx";
import {useApolloClient} from "@apollo/client";
import AuthService from "../auth/AuthService.tsx";
import LearningService from "../learning/LearningService.tsx";
import {Toast} from "primereact/toast";


const Navbar = () => {
	const [activeItem, setActiveItem] = useState('My Learning');
	const [showModal, setShowModal] = useState(false);
	const [username, setUsername] = useState('');
	const [challengeAvailable, setChallengeAvailable] = useState(false);
	const [role, setRole] = useState<number>(0);

	const client = useApolloClient();
	const authService = new AuthService(client);
	const learningService = new LearningService(client);

	const toast = useRef<Toast>(null);

	const {logout} = useAuth();

	const handleItemClick = (item: SetStateAction<string>) => {
		if (item === 'Daily Challenge' && !challengeAvailable) {
			toast.current?.show({
				severity: 'warn',
				summary: 'Daily Challenge',
				detail: 'Challenge not available',
				life: 3000
			});
			return;
		}

		if (item == 'Daily Challenge') {
			setShowModal(true);
		} else {
			setActiveItem(item);
		}
	};

	const closeModal = () => {
		setShowModal(false);
	};

	useEffect(() => {
		const token = authService.getUsername();

		if (token) {
			// @ts-ignore
			setUsername(token);
		} else {
			console.error('No token found');
		}

		const fetchChallengeAvailability = async () => {
			const data = await learningService.getNavbarData();
			setChallengeAvailable(data.dailyChallengeAvailable);
			setRole(data.role);
		};

		fetchChallengeAvailability();

	}, []);

	return (
		<nav className="sidebar">
			<Toast ref={toast}/>
			<DailyChallengeModal visible={showModal} onHide={closeModal}
								 onCorrectSubmit={() => setChallengeAvailable(false)}/>

			<div className="main-section">
				<div className="logo">
					<h2>Qriosity</h2>
				</div>
				<ul className="nav-list">
					<li className="nav-group">
						<h3>Main</h3>
						<ul>
							<li>
								<Link
									to="/home/my-learning"
									className={activeItem === 'My Learning' ? 'nav-link active' : 'nav-link '}
									onClick={() => handleItemClick('My Learning')}
								>
									<i className="pi pi-bookmark"></i> My Learning
								</Link>
							</li>
							<li>
								<a
									href="#"
									className={activeItem === 'Daily Challenge' ? 'nav-link active' : 'nav-link '}
									onClick={() => handleItemClick('Daily Challenge')}
									aria-disabled
								>
									<i className="pi pi-calendar-clock"></i> Daily Challenge
									{challengeAvailable && <span className="badge">New</span>}
								</a>
							</li>
						</ul>
					</li>
					<li className="nav-group">
						<h3>Explore</h3>
						<ul>
							<li>
								<Link
									to="/home/roadmaps"
									className={activeItem === 'Roadmaps' ? 'nav-link active' : 'nav-link'}
									onClick={() => handleItemClick('Roadmaps')}
								>
									<i className="pi pi-map"></i> Roadmaps
								</Link>
							</li>
							<li>
								<Link
									to="/home/courses"
									className={activeItem === 'Courses' ? 'nav-link active' : 'nav-link'}
									onClick={() => handleItemClick('Courses')}
								>
									<i className="pi pi-book"></i> Courses
								</Link>
							</li>
						</ul>
					</li>
					{role > 0 && (<li className="nav-group">
						<h3>Creator</h3>
						<ul>
							<li>
								<Link
									to="/home/build-roadmap"
									className={activeItem === 'Build Roadmap' ? 'nav-link active' : 'nav-link'}
									onClick={() => handleItemClick('Build Roadmap')}
								>
									<i className="pi pi-pencil"></i> Build Roadmap
								</Link>
							</li>
							<li>
								<Link
									to="/home/create-course"
									className={activeItem === 'Create Course' ? 'nav-link active' : 'nav-link'}
									onClick={() => handleItemClick('Create Course')}
								>
									<i className="pi pi-pen-to-square"></i> Create Course
								</Link>
							</li>
						</ul>
					</li>)}
				</ul>
			</div>
			<div className="user-section">
				<Link to={"/home/profile/" + authService.getCognitoUsername()} className="user-info">
					<img src="https://via.placeholder.com/40" alt="User Avatar" className="user-avatar"/>
					<span className="user-name">{username}</span>
				</Link>
				<button className="logout-btn" onClick={logout}>Logout</button>
			</div>
		</nav>
	);
};

export default Navbar;
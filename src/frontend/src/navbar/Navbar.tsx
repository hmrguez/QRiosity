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
	const [activeItem, setActiveItem] = useState('Dashboard');
	const [showModal, setShowModal] = useState(false);
	const [username, setUsername] = useState('');
	const [challengeAvailable, setChallengeAvailable] = useState(false);


	const client = useApolloClient();
	const authService = new AuthService(client);
	const learningService = new LearningService(client);


	const toast = useRef<Toast>(null);


	const {logout} = useAuth();

	const handleItemClick = (item: SetStateAction<string>) => {
		if (item === 'Daily Challenge' && !challengeAvailable) {
			toast.current?.show({severity: 'warn', summary: 'Daily Challenge', detail: 'Challenge not available', life: 3000});
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
			setUsername(token);
		} else {
			console.error('No token found');
		}

		const fetchChallengeAvailability = async () => {
			const available = await learningService.getDailyChallengeAvailability();
			setChallengeAvailable(available);
		};

		fetchChallengeAvailability();

	}, []);

	return (
		<div className="sidebar">
			<Toast ref={toast} />
			<nav className="mr-3">
				<div className="user-info flex flex-column mt-2">
					<span className="user-name text-center">{username}</span>
				</div>
				<ul className="nav-menu">
					<li className={activeItem === 'My Learning' ? 'active' : ''}
						onClick={() => handleItemClick('My Learning')}>
						<Link to="/my-learning"><i className="icon-dashboard"></i> My Learning</Link>
					</li>
				</ul>
				<h3 className="menu-header">LEARNING</h3>
				<ul className="nav-menu">
					<li className={activeItem === 'Learning Paths' ? 'active' : ''}
						onClick={() => handleItemClick('Learning Paths')}>
						<Link to="/learning-paths"><i className="icon-account"></i> Learning Paths</Link>
					</li>
					<li className={activeItem === 'Daily Challenge' ? 'active' : ''}
						onClick={() => handleItemClick('Daily Challenge')}>
						<a href="#" aria-disabled><i className="icon-publishing"></i> Daily
							Challenge {challengeAvailable &&
                                <span className="badge">New</span>}</a>
					</li>
				</ul>
				<h3 className="menu-header">OTHER MENU</h3>
				<ul className="nav-menu">
					<li className={activeItem === 'Logout' ? 'active' : ''} onClick={() => handleItemClick('Login')}>
						<a onClick={logout}><i className="icon-setting"></i>Logout</a>
					</li>
				</ul>
			</nav>
			<DailyChallengeModal visible={showModal} onHide={closeModal} onCorrectSubmit={() => setChallengeAvailable(false)}/>
		</div>
	);
};

export default Navbar;
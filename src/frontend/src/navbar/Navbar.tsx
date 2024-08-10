import {SetStateAction, useEffect, useState} from 'react';
import {Link} from 'react-router-dom';
import './Navbar.css';
import DailyChallengeModal from "../learning/DailyChallengeModal.tsx";
import {useAuth} from "../auth/AuthContext.tsx";
import {useApolloClient} from "@apollo/client";
import AuthService from "../auth/AuthService.tsx";

const Navbar = () => {
	const [activeItem, setActiveItem] = useState('Dashboard');
	const [showModal, setShowModal] = useState(false);
	const [username, setUsername] = useState('');

	const client = useApolloClient();
	const authService = new AuthService(client);


	const {logout} = useAuth();

	const handleItemClick = (item: SetStateAction<string>) => {
		if (item == 'Daily Challenge') {
			setShowModal(true);
		}
		setActiveItem(item);
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
	}, []);

	return (
		<div>

			<nav className="sidebar mr-3">
				<div className="user-info">
					<span className="user-name">{username}</span>
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
						<a href="#"><i className="icon-publishing"></i> Daily Challenge</a>
					</li>
				</ul>
				<h3 className="menu-header">OTHER MENU</h3>
				<ul className="nav-menu">
					<li className={activeItem === 'Logout' ? 'active' : ''} onClick={() => handleItemClick('Login')}>
						<a onClick={logout}><i className="icon-setting"></i>Logout</a>
					</li>
				</ul>
			</nav>
			<DailyChallengeModal visible={showModal} onHide={closeModal}/>
		</div>
	);
};

export default Navbar;
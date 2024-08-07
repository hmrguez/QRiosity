import {SetStateAction, useState} from 'react';
import {Link} from 'react-router-dom';
import './Navbar.css';
import DailyChallengeModal from "./DailyChallengeModal.tsx";

const Navbar = () => {
	const [activeItem, setActiveItem] = useState('Dashboard');
	const [showModal, setShowModal] = useState(false);

	const handleItemClick = (item: SetStateAction<string>) => {
		if (item == 'Daily Challenge') {
			setShowModal(true);
		}
		setActiveItem(item);
	};

	const closeModal = () => {
		setShowModal(false);
	};

	return (
		<div>

			<nav className="sidebar mr-3">
				<div className="user-info">
					<img src="assets/images.jpeg" alt="User Avatar" className="avatar"/>
					<span className="user-name">Wildan</span>
					<span className="user-role">Content Director</span>
				</div>
				<div className="search-bar">
					<input type="text" placeholder="Search"/>
				</div>
				<ul className="nav-menu">
					<li className={activeItem === 'Dashboard' ? 'active' : ''}
						onClick={() => handleItemClick('Dashboard')}>
						<Link to="/dashboard"><i className="icon-dashboard"></i> Dashboard</Link>
					</li>
					<li className={activeItem === 'Analytics' ? 'active' : ''}
						onClick={() => handleItemClick('Analytics')}>
						<Link to="/analytics"><i className="icon-analytics"></i> Analytics</Link>
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
					<li className={activeItem === 'Setting' ? 'active' : ''} onClick={() => handleItemClick('Setting')}>
						<Link to="/setting"><i className="icon-setting"></i> Setting</Link>
					</li>
				</ul>
			</nav>
			<DailyChallengeModal visible={showModal} onHide={closeModal}/>
		</div>
	);
};

export default Navbar;
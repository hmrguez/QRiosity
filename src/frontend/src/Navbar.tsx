import {SetStateAction, useState} from 'react';
import {Link} from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
	const [activeItem, setActiveItem] = useState('Dashboard');

	const handleItemClick = (item: SetStateAction<string>) => {
		setActiveItem(item);
	};

	return (
		<nav className="sidebar mr-3">
			<div className="user-info">
				<img src="user-avatar.jpg" alt="User Avatar" className="avatar"/>
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
			<h3 className="menu-header">ACCOUNT</h3>
			<ul className="nav-menu">
				<li className={activeItem === 'Account' ? 'active' : ''} onClick={() => handleItemClick('Account')}>
					<Link to="/account"><i className="icon-account"></i> Account</Link>
				</li>
				<li className={activeItem === 'My Publishing' ? 'active' : ''}
					onClick={() => handleItemClick('My Publishing')}>
					<Link to="/my-publishing"><i className="icon-publishing"></i> My Publishing</Link>
				</li>
				<li className={activeItem === 'Products' ? 'active' : ''}
					onClick={() => handleItemClick('Products')}>
					<Link to="/products"><i className="icon-products"></i> Products</Link>
				</li>
				<li className={activeItem === 'Orders' ? 'active' : ''} onClick={() => handleItemClick('Orders')}>
					<Link to="/orders"><i className="icon-orders"></i> Orders</Link>
				</li>
				<li className={activeItem === 'More' ? 'active' : ''} onClick={() => handleItemClick('More')}>
					<Link to="/more"><i className="icon-more"></i> More</Link>
				</li>
			</ul>
			<h3 className="menu-header">OTHER MENU</h3>
			<ul className="nav-menu">
				<li className={activeItem === 'Setting' ? 'active' : ''} onClick={() => handleItemClick('Setting')}>
					<Link to="/setting"><i className="icon-setting"></i> Setting</Link>
				</li>
				<li className={activeItem === 'Help' ? 'active' : ''} onClick={() => handleItemClick('Help')}>
					<Link to="/help"><i className="icon-help"></i> Help</Link>
				</li>
				<li className={activeItem === 'Subscriptions' ? 'active' : ''}
					onClick={() => handleItemClick('Subscriptions')}>
					<Link to="/subscriptions"><i className="icon-subscriptions"></i> Subscriptions</Link>
				</li>
			</ul>
		</nav>
	);
};

export default Navbar;
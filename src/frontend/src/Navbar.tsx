// src/frontend/src/components/Navbar.tsx
import React, {useState} from 'react';
import {Menubar} from 'primereact/menubar';
import {Badge} from 'primereact/badge';
import DailyChallengeModal from './DailyChallengeModal.tsx';


const Navbar: React.FC = () => {
	const [showModal, setShowModal] = useState(false);

	const items = [
		{
			label: 'Daily',
			icon: 'pi pi-calendar',
			command: () => setShowModal(true)
		}
	];

	return (
		<>
			<Menubar model={items} end={<Badge value=" " severity="danger"/>}/>
			<DailyChallengeModal visible={showModal} onHide={() => setShowModal(false)}/>
		</>
	);
};

export default Navbar;
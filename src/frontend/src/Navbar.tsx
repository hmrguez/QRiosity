import {useState} from 'react';
import {Menubar} from 'primereact/menubar';
import {Dialog} from 'primereact/dialog';

const Navbar = () => {
	const [visible, setVisible] = useState(false);

	const openModal = () => {
		setVisible(true);
	};

	const closeModal = () => {
		setVisible(false);
	};

	const items = [
		{
			label: 'Daily Challenge',
			icon: 'pi pi-fw pi-check',
			command: openModal
		}
	];

	return (
		<div>
			<Menubar model={items}/>
			<Dialog header="Daily Challenge" visible={visible} style={{width: '50vw'}} onHide={closeModal}>
				<p>Here is your daily challenge!</p>
			</Dialog>
		</div>
	);
};

export default Navbar;
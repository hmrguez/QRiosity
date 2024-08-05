// src/components/Navbar.tsx
import { useState } from 'react';
import { Menubar } from 'primereact/menubar';
import DailyChallengeModal from './DailyChallenge';

const Navbar = () => {
  const [visible, setVisible] = useState<boolean>(false);

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
      <Menubar model={items} />
      <DailyChallengeModal visible={visible} onHide={closeModal} />
    </div>
  );
};

export default Navbar;
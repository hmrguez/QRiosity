// src/components/DailyChallenge.js
import { useState } from 'react';
import { Button } from 'primereact/button';
import { Dialog } from 'primereact/dialog';

const DailyChallenge = () => {
  const [visible, setVisible] = useState(false);

  const openModal = () => {
    setVisible(true);
  };

  const closeModal = () => {
    setVisible(false);
  };

  return (
    <div>
      <Button label="Daily Challenge" icon="pi pi-check" onClick={openModal} />
      <Dialog header="Daily Challenge" visible={visible} style={{ width: '50vw' }} onHide={closeModal}>
        <p>Here is your daily challenge!</p>
      </Dialog>
    </div>
  );
};

export default DailyChallenge;
// src/frontend/src/components/DailyChallengeModal.tsx
import React, { useState } from 'react';
import { Dialog } from 'primereact/dialog';
import DailyChallenge from './DailyChallenge';

interface DailyChallengeModalProps {
    visible: boolean;
    onHide: () => void;
}

const DailyChallengeModal: React.FC<DailyChallengeModalProps> = ({ visible, onHide }) => {
    const [answer, setAnswer] = useState('');

    return (
        <Dialog header="Daily Challenge" visible={visible} style={{ width: '50vw' }} onHide={onHide}>
            <DailyChallenge answer={answer} setAnswer={setAnswer} />
        </Dialog>
    );
};

export default DailyChallengeModal;
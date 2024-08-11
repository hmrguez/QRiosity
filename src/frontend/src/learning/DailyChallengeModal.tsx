// src/frontend/src/components/DailyChallengeModal.tsx
import React, { useState } from 'react';
import { Dialog } from 'primereact/dialog';
import DailyChallenge from './DailyChallenge';
import "./DailyChallenge.css"

interface DailyChallengeModalProps {
    visible: boolean;
    onHide: () => void;
    onCorrectSubmit: () => void;
}

const DailyChallengeModal: React.FC<DailyChallengeModalProps> = ({ visible, onHide, onCorrectSubmit }) => {
    const [answer, setAnswer] = useState('');

    return (
        <Dialog visible={visible} style={{ width: '50vw' }} onHide={onHide}>
            <DailyChallenge answer={answer} setAnswer={setAnswer} onHide={onHide} onCorrectSubmit={onCorrectSubmit} />
        </Dialog>
    );
};

export default DailyChallengeModal;
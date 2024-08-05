import {FC} from 'react';
import {Dialog} from 'primereact/dialog';

interface DailyChallengeModalProps {
	visible: boolean;
	onHide: () => void;
}

const DailyChallengeModal: FC<DailyChallengeModalProps> = ({visible, onHide}) => {
	return (
		<Dialog header="Daily Challenge" visible={visible} style={{width: '50vw'}} onHide={onHide}>
			<p>Here is your daily challenge!</p>
		</Dialog>
	);
};

export default DailyChallengeModal;
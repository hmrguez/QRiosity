import {useRef} from 'react';
import FileUpload from "./FileInput.tsx";

const ParentComponent = () => {
	const fileUploadRef = useRef<{ handleSubmit: () => void }>(null);

	const callHandleSubmit = () => {
		if (fileUploadRef.current) {
			fileUploadRef.current.handleSubmit();
		}
	};

	return (
		<div>
			<h1>Parent Component</h1>
			<button onClick={callHandleSubmit}>Submit from Parent</button>
			<FileUpload ref={fileUploadRef}/>
		</div>
	);
};

export default ParentComponent;
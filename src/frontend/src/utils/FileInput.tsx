import React, {forwardRef, useImperativeHandle, useState} from 'react';
import './FileInput.css';

const FileUpload = forwardRef((_, ref) => {
	const [selectedFile, setSelectedFile] = useState<File | null>(null);
	const [filename, setFilename] = useState<string>('');

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setFilename(event.target.files ? event.target.files[0].name : '');
		setSelectedFile(event.target.files ? event.target.files[0] : null);
	};

	const handleSubmit = async (event?: React.FormEvent): Promise<string | null> => {
		if (event) event.preventDefault();

		if (!selectedFile) {
			alert("Please select a file to upload.");
			return null;
		}

		const formData = new FormData();
		formData.append('file', selectedFile);

		try {
			const response = await fetch('https://kqssn9e4a3.execute-api.us-east-2.amazonaws.com/image-upload', {
				method: 'POST',
				body: formData,
			});

			const data = await response.json();
			console.log('File uploaded successfully:', data);

			return data.filename;

		} catch (error) {
			console.error('Error uploading file:', error);
			return null;
		}
	};

	useImperativeHandle(ref, () => ({
		handleSubmit,
	}));

	return (

		<form onSubmit={handleSubmit} className="file-input-container">
			<input type="file" id="file-input" className="file-input" onChange={handleFileChange}/>
			<label htmlFor="file-input" className="file-input-label">Thumbnail</label>
			<div id="file-name" className="file-name">{filename}</div>
		</form>

	)
});

export default FileUpload;
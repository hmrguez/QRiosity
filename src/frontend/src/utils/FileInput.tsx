import React, {forwardRef, useImperativeHandle, useState} from 'react';
import './FileInput.css';

const FileUpload = forwardRef((_, ref) => {
	const [selectedFile, setSelectedFile] = useState<File | null>(null);
	const [filename, setFilename] = useState<string>('');

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files ? event.target.files[0] : null;
		if (file) {
			const fileType = file.type.toLowerCase();
			const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];

			if (!validTypes.includes(fileType)) {
				alert('Invalid file type. Please select a JPEG, JPG, PNG, or WEBP image.');
				return;
			}

			if (file.size > 5 * 1024 * 1024) {
				alert('File size exceeds the 5MB limit.');
				return;
			}

			setFilename(file.name);
			setSelectedFile(file);
		} else {
			setFilename('');
			setSelectedFile(null);
		}
	};

	const handleSubmit = () => {
		return new Promise((resolve, reject) => {
			if (!selectedFile) {
				reject(new Error('No file selected'));
				return;
			}

			// Check if the file type is supported
			const supportedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp'];
			if (!supportedTypes.includes(selectedFile.type)) {
				reject(new Error('Unsupported file type'));
				return;
			}

			const reader = new FileReader();
			reader.readAsDataURL(selectedFile);
			reader.onloadend = async () => {
				const base64data = (reader.result as string).split(',')[1];

				try {
					const response = await fetch('https://kqssn9e4a3.execute-api.us-east-2.amazonaws.com/roadmap-image', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({
							image: base64data,
							mimeType: selectedFile.type
						}),
					});

					if (!response.ok) {
						throw new Error(`HTTP error! status: ${response.status}`);
					}

					const data = await response.json();

					console.log('Image uploaded successfully:', data.url);

					resolve(data.url);
				} catch (error: any) {
					console.log('Error uploading image: ' + error.message);
					reject(error);
				}
			};
		});
	};

	// const handleSubmit = async (event?: React.FormEvent): Promise<string | null> => {
	//
	// 	if (event) event.preventDefault();
	//
	//
	// 	if (!selectedFile) {
	// 		return null;
	// 	}
	//
	// 	const reader = new FileReader();
	// 	reader.readAsDataURL(selectedFile);
	//
	// 	return new Promise((resolve, reject) => {
	// 		reader.onloadend = async () => {
	// 			const base64Data = (reader.result as string).split(',')[1]; // get base64 data without the data URL prefix
	// 			const contentType = selectedFile.type;
	//
	// 			try {
	// 				const response = await fetch('https://kqssn9e4a3.execute-api.us-east-2.amazonaws.com/roadmap-image', {
	// 					method: 'POST',
	// 					headers: {
	// 						'Content-Type': 'application/json',
	// 					},
	// 					body: JSON.stringify({
	// 						body: base64Data,
	// 						isBase64Encoded: true,
	// 						headers: {'Content-Type': contentType},
	// 					}),
	// 				});
	//
	// 				const result = await response.text();
	// 				if (response.ok) {
	// 					console.log('File uploaded successfully:', result);
	// 					resolve(result);
	// 				} else {
	// 					console.log('Error uploading file:', result);
	// 					resolve(null);
	// 				}
	// 			} catch (error) {
	// 				console.error('Error uploading file:', error);
	// 				reject(error);
	// 			}
	// 		};
	// 	});
	// };

	// const handleSubmit = async (event?: React.FormEvent): Promise<string | null> => {
	// 	if (event) event.preventDefault();
	//
	// 	if (!selectedFile) {
	// 		alert("Please select a file to upload.");
	// 		return null;
	// 	}
	//
	// 	const reader = new FileReader();
	// 	reader.readAsDataURL(selectedFile);
	//
	// 	// const formData = new FormData();
	// 	// formData.append('file', selectedFile);
	//
	// 	try {
	//
	// 		console.log(selectedFile.type)
	//
	// 		const response = await fetch('https://kqssn9e4a3.execute-api.us-east-2.amazonaws.com/roadmap-image', {
	// 			method: 'POST',
	// 			body: formData,
	// 			headers: {
	// 				'Content-Type': selectedFile.type,
	// 			},
	// 		});
	//
	// 		if (!response.ok) {
	// 			throw new Error(`HTTP error! status: ${response.status}`);
	// 		}
	//
	// 		const data = await response.text();
	// 		console.log('File uploaded successfully:', data);
	//
	// 		return data;
	// 	} catch (error) {
	// 		console.error('Error uploading file:', error);
	// 		return null;
	// 	}
	// };

	useImperativeHandle(ref, () => ({
		handleSubmit,
	}));

	return (
		<form onSubmit={handleSubmit} className="file-input-container">
			<input
				type="file"
				id="file-input"
				className="file-input"
				onChange={handleFileChange}
				accept=".jpeg,.jpg,.png,.webp"
			/>
			<label htmlFor="file-input" className="file-input-label">Thumbnail</label>
			<div id="file-name" className="file-name">{filename}</div>
		</form>
	);
});

export default FileUpload;
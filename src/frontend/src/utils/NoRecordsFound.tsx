import NoRecordsImage from '../assets/not-found.png';
import "./NoRecordsFound.css"

const NoRecordsFound = () => {
	return (
		<div className="not-found-container">
			<h4>No records found</h4>
			<img className="not-found-image" src={NoRecordsImage} alt="No Records Found"/>
		</div>
	);
};


export default NoRecordsFound;
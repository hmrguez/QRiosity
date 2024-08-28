import "./MyLearning.css";
import RoadmapList from "./RoadmapList";

const MyLearning = () => {
	const roadmaps = [
		{
			id: "1",
			title: "Web Development Fundamentals",
			author: "John Doe",
			topics: ["HTML", "CSS", "JavaScript"],
			difficulty: "Beginner",
			likes: 150,
		},
		{
			id: "2",
			title: "Advanced Machine Learning",
			author: "Jane Smith",
			topics: ["Python", "TensorFlow", "Neural Networks"],
			difficulty: "Advanced",
			likes: 320,
			isCustom: true,
		},
	];

	return (
		<div className="my-learning">
			<RoadmapList roadmaps={roadmaps}/>
		</div>
	);
};

export default MyLearning;
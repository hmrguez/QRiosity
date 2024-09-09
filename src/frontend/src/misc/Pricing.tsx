// Define the pricing plans data
import {Tag} from "primereact/tag";
import './Pricing.css';

const pricingPlans = [
	{
		title: "Student",
		price: "Free",
		features: ["Access to 5 roadmaps/month", "1 daily challenge", "1 roadmap progress tracking at a time"],
		buttonText: "Select Plan"
	},
	{
		title: "Apprentice",
		price: "$3.99",
		period: "/month",
		features: ["Access to 25 roadmaps/month", "5 daily challenges", "5 roadmaps progress tracking at a time", "Create 2 roadmaps / month", "3 autogen roadmap tickets / month"],
		buttonText: "Select Plan"
	},
	{
		title: "Creator",
		price: "$8.99",
		period: "/month",
		features: ["Everything on Apprentice", "Unlimited liked roadmaps", "10 autogen roadmap tickets / month", "Exclusive verified badge", "Enhanced analytics"],
		buttonText: "Select Plan",
		popular: true,
		beta: true
	},
	// {
	// 	title: "Enterprise",
	// 	price: "Custom Pricing",
	// 	features: ["Access for teams", "Dedicated support", "Custom daily challenges"],
	// 	buttonText: "Contact Us"
	// }
];

const Pricing = () => {
	return (
		<section className="pricing-section">
			<h1 className="pricing-title">Choose Your Plan</h1>
			<div className="pricing-container">
				{pricingPlans.map((plan, index) => (
					<div key={index} className={`pricing-card`}>
						<h2 className="plan-title">{plan.title}</h2>
						{plan.beta &&
                            <Tag style={{background: 'gold'}}>Beta</Tag>} {/* Conditionally render beta tag */}
						<p className="price">{plan.price}<span>{plan.period}</span></p>
						<ul className="features">
							{plan.features.map((feature, idx) => (
								<li key={idx}>{feature}</li>
							))}
						</ul>
						<button className="select-button">{plan.buttonText}</button>
					</div>
				))}
			</div>
			<footer className="footer">
				<p>&copy; 2024 Qriosity. All rights reserved.</p>
				<ul className="footer-links">
					{/*<li><a href="/terms">Terms of Service</a></li>*/}
					{/*<li><a href="/privacy">Privacy Policy</a></li>*/}
					<li><a href="/">About Us</a></li>
				</ul>
			</footer>
		</section>
	);
};

export default Pricing;
import './LandingPage.css';
import image from './assets/image.jpeg'
import {Link} from "react-router-dom";

const LandingPage = () => {
	return (
		<div className="landing">
			<header>
				<div className="container">
					<nav>
						<div className="logo">QRIOSITY</div>
						<div className="nav-links">
							<a className="nav-link" href="#about">About</a>
							<a className="nav-link" href="#features">Features</a>
							<a className="nav-link" href="#contact">Contact</a>

							<div className="auth-buttons">
								<Link to="/login" className="auth-button">Login</Link>
								<Link to="/register" className="auth-button">Signup</Link>
							</div>
						</div>
					</nav>
				</div>
			</header>
			<main className="container">
				<section id="about" className="hero">
					<div className="hero-content">
						<h1>Empower Your Learning Journey with Qriosity</h1>
						<p>Unlock personalized learning experiences with AI-driven insights and custom paths.</p>
						<Link to="/register" className="cta-button">Start Your Journey</Link>
					</div>
					<div className="hero-image">
						<img src={image} alt="AI-powered learning illustration" />
					</div>
				</section>
				<section id="features" className="features">
					<div className="feature">
						<h2>Custom Learning Paths</h2>
						<p>Create and follow learning paths tailored to your goals and background.</p>
					</div>
					<div className="feature">
						<h2>Daily Challenges</h2>
						<p>Engage with AI-generated challenges to test your knowledge and track progress.</p>
					</div>
					<div className="feature">
						<h2>Extensive Content Library</h2>
						<p>Explore a vast library of courses and resources from trusted providers.</p>
					</div>
				</section>
				<section className="contact" id="contact">
					<h2>Contact Us</h2>
					<p>Have questions or need more information? Reach out to us, and we’ll be happy to assist you.</p>
					<form>
						<input type="text" name="name" placeholder="Your Name" required />
						<input type="email" name="email" placeholder="Your Email" required />
						<textarea name="message" placeholder="Your Message" required></textarea>
						<input type="submit" value="Send Message" />
					</form>
				</section>
			</main>
			<footer>
				<div className="container">
					<p>&copy; 2024 Qriosity. All rights reserved.</p>
				</div>
			</footer>
		</div>
	);
}

export default LandingPage;
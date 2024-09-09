import "./OutOfViews.css"

const OutOfViews = () => {
	return (
		<div className="out-card">
			<p className="out-card-text">No roadmaps views remaining, consider upgrading your plan</p>
			<a href="/pricing" className="out-subscribe-button">Subscribe</a>
		</div>
	)
}

export default OutOfViews;
import './App.css'
import {Outlet} from "react-router-dom";
import Navbar from "./navbar/Navbar.tsx";

function App() {
	return (
		<div className="flex gap-3">
			<Navbar/>
			<Outlet/>
		</div>
	);
}

export default App;
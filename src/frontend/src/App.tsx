import './App.css'
import {Outlet} from "react-router-dom";
import Navbar from "./navbar/Navbar.tsx";

function App() {
	return (
		<div className="flex gap-3">
			<Navbar/>
			{/*<div style={{marginLeft: '275px', marginRight: "30px", marginTop: "-17px"}}>*/}
			<Outlet/>
			{/*</div>*/}
		</div>
	);
}

export default App;
import { Link } from "react-router-dom";
import "../styles/Navbar.css";
function Navbar() {
  return (
    <div className="navbar">
      <Link to="/" className="logo">
        Inventory System
      </Link>
      <div className="nav-items">
        <Link to="/" className="nav-item">
          Products
        </Link>
        <Link to="/category" className="nav-item">
          Categories
        </Link>
        <Link to="/reports" className="nav-item">
          Reports
          </Link>
      </div>
    </div>
  );
}

export default Navbar;

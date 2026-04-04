import { Link ,  } from "react-router-dom";
import type { Category } from "../interfaces/inventory";
import '../styles/Card.css'
function CategoryCard({id,name ,description } : Category) {


  return (
    <div className="card" >
      <p className="name">Name : {name}</p>
      <p className="description">Description : {description}</p>
      <Link to={`/category/${id}`}><button className="view-category">View Category</button></Link>
      <Link to={`/?category=${id}`}><button className="view-product">View Products</button></Link>
    </div>
  )
}

export default CategoryCard

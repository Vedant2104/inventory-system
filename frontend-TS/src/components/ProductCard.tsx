import { Link } from "react-router-dom";
import type { Product } from "../interfaces/inventory";
import '../styles/Card.css'
import { useState } from "react";
function ProductCard({id,name ,description , category , price, brand , quantity} : Product) {

  const [expanded, setExpanded] = useState<boolean>(false);

  return (
    <div className="card">
      <p className="name">Name : {name}</p>
      <p className="description">Description : {description}</p>

      <div className={`details ${expanded ? "expanded" : ""}`}>
        <p className="category" data-tooltip = {category.description}>Category : {category.name}</p>
        <p className="price">Price : {price}</p>
        <p className="brand ">Brand: {brand}</p>
        <p className="quantity">Quantity : {quantity}</p>
      </div>
      
      <button className="read-more" onClick={()=>{setExpanded(!expanded)}}>{expanded ? "Read Less" : "Read More"}</button>  
      <Link to={`/product/${id}`}><button className="view-product">View Product</button></Link>
    </div>
  )
}

export default ProductCard

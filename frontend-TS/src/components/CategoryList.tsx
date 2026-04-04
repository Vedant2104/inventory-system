import { useEffect, useState } from "react";
import type { Category } from "../interfaces/inventory";
import Loader from "./Loader";
import CategoryCard from "./CategoryCard";
import "../styles/Product_list.css"
import { useNavigate } from "react-router-dom";
export default function CategoryList() {
  const [categories , setCategories] = useState<Category[]>([])
  const [Loading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");
  const navigate = useNavigate();
  function handleFetchCategories(){
    setIsLoading(true);

    fetch("http://localhost:8080/category", { method: "GET" })
    .then((res) => {
      if(!res.ok){
        throw new Error("Failed to fetch categories");
      }
      return res.json();
    })
    .then(data=>{
      setCategories(data || []);
    })
    .catch(err=>{
      console.log(err);
      setError(err.message);
    })
    .finally(()=>{
      setIsLoading(false);
    })
  }

  useEffect(()=>{handleFetchCategories()},[])

  if(Loading){
    return <Loader />;
  }
  if(error){
    return <h1 className="error">{error}</h1>
  }

  return (
    <div className="product-list">
      <button className="add-item" onClick={() => navigate("/category/new")}>➕Add Category</button>
      {categories.map((category) => (
        <CategoryCard key={category.id} {...category} />
      ))}
    </div>
  )
}
